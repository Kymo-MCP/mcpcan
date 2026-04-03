package container

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/rest"

	"github.com/kymo-mcp/mcpcan/pkg/k8s"
)

// KubernetesRuntime Kubernetes runtime implementation
type KubernetesRuntime struct {
	Entry *k8s.Entry
}

// NewKubernetesRuntime creates Kubernetes runtime
func NewKubernetesRuntime(kubeconfig *rest.Config, namespace string) (*KubernetesRuntime, error) {
	k8sEntry, err := k8s.NewEntry(kubeconfig, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes client: %w", err)
	}
	return &KubernetesRuntime{
		Entry: k8sEntry,
	}, nil
}

// GetContainerManager gets container manager
func (kr *KubernetesRuntime) GetContainerManager() ContainerManager {
	return &KubernetesContainerManager{Entry: kr.Entry}
}

// GetServiceManager gets service manager
func (kr *KubernetesRuntime) GetServiceManager() ServiceManager {
	return &KubernetesServiceManager{Entry: kr.Entry}
}

// GetVolumeManager gets volume manager
func (kr *KubernetesRuntime) GetVolumeManager() VolumeManager {
	return &KubernetesVolumeManager{Entry: kr.Entry}
}

// GetRuntimeType gets runtime type
func (kr *KubernetesRuntime) GetRuntimeType() ContainerRuntime {
	return RuntimeKubernetes
}

// KubernetesContainerManager Kubernetes container manager implementation
type KubernetesContainerManager struct {
	Entry *k8s.Entry
}

// sanitizeK8sLabels strips labels that are not compatible with Kubernetes metadata constraints.
// Key points:
// 1) Docker-only Traefik labels (traefik.*) are intentionally dropped in K8s mode.
// 2) Any key/value that fails K8s label validation is dropped to prevent Deployment creation failures.
func sanitizeK8sLabels(labels map[string]string) map[string]string {
	if len(labels) == 0 {
		return nil
	}

	sanitized := make(map[string]string, len(labels))
	for k, v := range labels {
		// Docker provider labels are not valid/needed in Kubernetes mode.
		if strings.HasPrefix(k, "traefik.") {
			continue
		}
		if errs := validation.IsQualifiedName(k); len(errs) > 0 {
			continue
		}
		if errs := validation.IsValidLabelValue(v); len(errs) > 0 {
			continue
		}
		sanitized[k] = v
	}
	return sanitized
}

// Create creates container (Deployment)
func (kcm *KubernetesContainerManager) Create(ctx context.Context, options ContainerCreateOptions) (string, error) {
	// Initialize basic DeploymentCreateOptions
	deploymentOptions := k8s.DeploymentCreateOptions{
		ImageName: options.ImageName,
		AppName:   options.ContainerName,
		Namespace: kcm.Entry.Namespace,
		Port:      options.Port,
		Replicas:  1, // default single replica
	}

	// Set execution command (if specified)
	if len(options.Command) > 0 {
		deploymentOptions.Command = options.Command
	}

	// Set command arguments (if specified)
	if len(options.CommandArgs) > 0 {
		deploymentOptions.CommandArgs = options.CommandArgs
	}

	// Set environment variables (if exist)
	if len(options.EnvVars) > 0 {
		deploymentOptions.EnvVars = options.EnvVars
	}

	// Set working directory (if specified)
	if options.WorkingDir != "" {
		deploymentOptions.WorkingDir = options.WorkingDir
	}

	// Set labels (if exist)
	if len(options.Labels) > 0 {
		deploymentOptions.Labels = sanitizeK8sLabels(options.Labels)
	}

	// Set annotations (if exist) — K8s 下 Traefik 路由配置通过 annotations 传递，规避 label 校验限制
	if len(options.Annotations) > 0 {
		deploymentOptions.Annotations = options.Annotations
	}

	// Set restart policy (if specified)
	if options.RestartPolicy != "" {
		deploymentOptions.RestartPolicy = options.RestartPolicy
	}

	// Handle volume mounts and affinity configuration
	if len(options.Mounts) > 0 {
		deploymentOptions.VolumeMounts = options.Mounts
	}

	// Set readiness probe
	if options.ReadinessProbe != nil {
		deploymentOptions.ReadinessProbe = options.ReadinessProbe
	}

	// Set image pull secrets
	if len(options.ImagePullSecrets) > 0 {
		deploymentOptions.ImagePullSecrets = options.ImagePullSecrets
	}

	// Set sidecar config
	if options.Sidecar != nil {
		// In Kubernetes, the sidecar and main container are deployed within the same Pod.
		// Map the target URL to use localhost instead of querying DNS for the container hostname.
		if targetUrl, ok := options.Sidecar.EnvVars["MCP_TARGET_URL"]; ok {
			options.Sidecar.EnvVars["MCP_TARGET_URL"] = strings.Replace(targetUrl, options.ContainerName, "localhost", 1)
		}

		deploymentOptions.Sidecar = &k8s.SidecarOptions{
			ImageName:     options.Sidecar.ImageName,
			ContainerName: options.Sidecar.ContainerName,
			Port:          options.Sidecar.Port,
			Command:       options.Sidecar.Command,
			CommandArgs:   options.Sidecar.CommandArgs,
			EnvVars:       options.Sidecar.EnvVars,
		}
	}

	// Create deployment
	deploymentName, err := kcm.Entry.Client.Deployment().Create(deploymentOptions)
	if err != nil {
		return "", fmt.Errorf("failed to create Deployment: %w", err)
	}
	return deploymentName, nil
}

// Delete deletes container (Deployment)
func (kcm *KubernetesContainerManager) Delete(ctx context.Context, containerName string) error {
	return kcm.Entry.Client.Deployment().Delete(containerName)
}

// Scale sets container replica count (Deployment)
func (kcm *KubernetesContainerManager) Scale(ctx context.Context, containerName string, replicas int32) error {
	return kcm.Entry.Client.Deployment().Scale(containerName, replicas)
}

// Restart restarts container (K8s environment: delete and recreate if exists, create directly if not exists)
func (kcm *KubernetesContainerManager) Restart(ctx context.Context, options ContainerCreateOptions) error {
	// Check if deployment exists
	_, err := kcm.Entry.Client.Deployment().Get(options.ContainerName)
	if err == nil {
		// deployment exists, delete first
		if err := kcm.Delete(ctx, options.ContainerName); err != nil {
			return fmt.Errorf("failed to delete existing deployment: %w", err)
		}

		// Start async process to probe and create
		go kcm.asyncProbeAndCreate(context.Background(), options)
		return nil
	} else if !isNotFoundError(err) {
		// If not NotFound error, it might be other issues (like network problems, permission issues, etc.)
		return fmt.Errorf("failed to check deployment status: %w", err)
	}

	// deployment doesn't exist, create directly
	_, createErr := kcm.Create(ctx, options)
	if createErr != nil {
		return fmt.Errorf("failed to create deployment: %w", createErr)
	}

	return nil
}

// asyncProbeAndCreate probes for deployment deletion and creates new deployment
func (kcm *KubernetesContainerManager) asyncProbeAndCreate(ctx context.Context, options ContainerCreateOptions) {
	const (
		maxProbes     = 5               // maximum probe count
		probeInterval = 5 * time.Second // probe interval
	)

	for i := 0; i < maxProbes; i++ {
		// Wait for probe interval
		time.Sleep(probeInterval)

		// Check if deployment still exists
		deployment, err := kcm.Entry.Client.Deployment().Get(options.ContainerName)
		if err != nil {
			// If get fails and is NotFound error, deletion is successful
			if isNotFoundError(err) {
				// Deployment doesn't exist, create new one
				if _, createErr := kcm.Create(ctx, options); createErr != nil {
					// Log error but don't return, as this is async
					fmt.Printf("Failed to create deployment in async process: %v\n", createErr)
				}
				return
			}
			// Other errors, continue probing
			continue
		}

		// Check if deployment is being deleted (has DeletionTimestamp)
		if deployment.DeletionTimestamp != nil {
			// Being deleted, continue probing
			continue
		}
	}

	// If we reach here, probing timed out
	fmt.Printf("Async probe timed out after %d attempts for deployment %s\n", maxProbes, options.ContainerName)
}

// waitForDeploymentDeletion waits for deployment to be completely deleted
func (kcm *KubernetesContainerManager) waitForDeploymentDeletion(ctx context.Context, deploymentName string, resourceVersion string) error {
	const (
		maxRetries    = 30              // maximum retry count
		retryInterval = 1 * time.Second // retry interval
	)

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for deployment deletion was cancelled: %w", ctx.Err())
		default:
		}

		// Check if deployment still exists
		deployment, err := kcm.Entry.Client.Deployment().Get(deploymentName)
		if err != nil {
			// If get fails and is NotFound error, deletion is successful
			if isNotFoundError(err) {
				return nil
			}
			// Other errors continue retrying
		} else {
			// If deployment still exists, check if it's a new instance (by ResourceVersion)
			if deployment.ResourceVersion != resourceVersion {
				// Different ResourceVersion means it's a newly created deployment, deletion is complete
				return nil
			}

			// Check if deployment is being deleted (has DeletionTimestamp)
			if deployment.DeletionTimestamp != nil {
				// Being deleted, continue waiting
				time.Sleep(retryInterval)
				continue
			}

			// deployment still exists and has no deletion mark, deletion might have failed
		}

		time.Sleep(retryInterval)
	}
	return fmt.Errorf("waiting for deployment deletion timed out, exceeded %d seconds", maxRetries)
}

// GetInfo gets container information
func (kcm *KubernetesContainerManager) GetInfo(ctx context.Context, containerName string) (*ContainerInfo, error) {
	deployment, err := kcm.Entry.Client.Deployment().Get(containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get Deployment information: %w", err)
	}

	pods, err := kcm.Entry.Client.Deployment().GetPods(deployment.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pod list: %w", err)
	}

	status := "Unknown"
	ip := ""

	if len(pods) > 0 {
		// Take the first Pod as status reference
		pod := pods[0]
		status = string(pod.Status.Phase)
		ip = pod.Status.PodIP
	}

	// Extract ports
	var ports []int32
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		for _, p := range deployment.Spec.Template.Spec.Containers[0].Ports {
			ports = append(ports, p.ContainerPort)
		}
	}

	return &ContainerInfo{
		Name:      deployment.Name,
		Status:    status,
		IP:        ip,
		Ports:     ports,
		Labels:    deployment.Labels,
		CreatedAt: deployment.CreationTimestamp.Format(time.RFC3339),
	}, nil
}

// IsReady checks if container is ready
// IsReady checks if container is ready
func (kcm *KubernetesContainerManager) IsReady(ctx context.Context, containerName string) (bool, string, error) {

	deployment, err := kcm.Entry.Client.Deployment().Get(containerName)
	if err != nil {
		return false, "", err
	}

	if deployment.Status.ReadyReplicas == deployment.Status.Replicas && deployment.Status.Replicas > 0 {
		return true, "Running", nil
	}

	return false, fmt.Sprintf("Ready: %d/%d", deployment.Status.ReadyReplicas, deployment.Status.Replicas), nil
}

// GetEvents gets container events
func (kcm *KubernetesContainerManager) GetEvents(ctx context.Context, containerName string) ([]ContainerEvent, error) {
	deployment, err := kcm.Entry.Client.Deployment().Get(containerName)
	if err != nil {
		return nil, err
	}
	pods, err := kcm.Entry.Client.Deployment().GetPods(deployment.Name)
	if err != nil {
		return nil, err
	}

	if len(pods) == 0 {
		return []ContainerEvent{}, nil
	}

	k8sEvents, err := kcm.Entry.Pod.GetEvents(pods[0].Name, string(pods[0].UID))
	if err != nil {
		return nil, err
	}

	var events []ContainerEvent
	for _, e := range k8sEvents {
		events = append(events, ContainerEvent{
			Type:      e.Type,
			Reason:    e.Reason,
			Message:   e.Message,
			Timestamp: e.Timestamp.Unix(),
		})
	}

	return events, nil
}

// GetLogs gets container logs
// K8S 实例 Pod 中包含主容器与 sidecar 容器，必须显式指定主容器名，否则 K8S API 会因歧义报错。
// containerName 格式如 "mcp-instance-xxx-container"，即主容器名（不含 -sidecar 后缀）。
func (kcm *KubernetesContainerManager) GetLogs(ctx context.Context, containerName string, lines int64) (string, error) {
	deployment, err := kcm.Entry.Client.Deployment().Get(containerName)
	if err != nil {
		return "", err
	}
	pods, err := kcm.Entry.Client.Deployment().GetPods(deployment.Name)
	if err != nil {
		return "", err
	}

	if len(pods) == 0 {
		return "", fmt.Errorf("no pods found for deployment %s", containerName)
	}

	// 多容器 Pod 需显式指定主容器名（containerName 即主容器名，不含 -sidecar）
	return kcm.Entry.Pod.GetLogs(pods[0].Name, lines, containerName)
}

// GetWarningEvents gets container warning events
func (kcm *KubernetesContainerManager) GetWarningEvents(ctx context.Context, containerName string) ([]ContainerEvent, error) {
	events, err := kcm.GetEvents(ctx, containerName)
	if err != nil {
		return nil, err
	}

	var warnings []ContainerEvent
	for _, e := range events {
		if e.Type == "Warning" {
			warnings = append(warnings, e)
		}
	}

	return warnings, nil
}

// ListByLabel lists all Deployments in the namespace that match ALL of the given labels.
func (kcm *KubernetesContainerManager) ListByLabel(ctx context.Context, labels map[string]string) ([]ContainerInfo, error) {
	// Build label selector string, e.g. "managed-by=mcpcan"
	var parts []string
	for k, v := range labels {
		parts = append(parts, k+"="+v)
	}
	selector := strings.Join(parts, ",")

	deployments, err := kcm.Entry.Client.Deployment().ListByLabelSelector(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to list Deployments by label: %w", err)
	}

	var result []ContainerInfo
	for _, d := range deployments {
		result = append(result, ContainerInfo{
			Name:      d.Name,
			Status:    "unknown",
			Labels:    d.Labels,
			CreatedAt: d.CreationTimestamp.Format(time.RFC3339),
		})
	}
	return result, nil
}

// KubernetesServiceManager Kubernetes service manager implementation
type KubernetesServiceManager struct {
	Entry *k8s.Entry
}

// Create 创建 Service，若同名 Service 已存在则直接返回（幂等）。
// 编辑实例时会重走 CreateContainer 流程，同名 Service 已存在属正常情况，不应报错。
func (ksm *KubernetesServiceManager) Create(ctx context.Context, serviceName string, port int32, selector map[string]string) (*ServiceInfo, error) {
	cleanSelector := sanitizeK8sLabels(selector)
	if len(cleanSelector) == 0 {
		return nil, fmt.Errorf("failed to create Service: selector is empty after sanitization")
	}

	svcCfg := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None", // Headless Service
			Ports: []corev1.ServicePort{
				{
					Port: port,
				},
			},
			Selector: cleanSelector,
		},
	}

	service, err := ksm.Entry.Service.Create(svcCfg)
	if err != nil {
		// 若 Service 已存在，直接 Get 返回，避免幂等场景下误报错误
		if strings.Contains(err.Error(), "already exists") {
			existing, getErr := ksm.Get(ctx, serviceName)
			if getErr == nil {
				return existing, nil
			}
		}
		return nil, fmt.Errorf("failed to create Service: %w", err)
	}

	// Extract port information
	var ports []int32
	for _, port := range service.Spec.Ports {
		ports = append(ports, port.Port)
	}

	return &ServiceInfo{
		Name:      service.Name,
		ClusterIP: service.Spec.ClusterIP,
		Ports:     ports,
		Labels:    service.Labels,
	}, nil
}


// Delete deletes service
func (ksm *KubernetesServiceManager) Delete(ctx context.Context, serviceName string) error {
	return ksm.Entry.Service.Delete(serviceName)
}

// Get gets service information
func (ksm *KubernetesServiceManager) Get(ctx context.Context, serviceName string) (*ServiceInfo, error) {
	service, err := ksm.Entry.Service.Get(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get Service information: %w", err)
	}

	// Extract port information
	var ports []int32
	for _, port := range service.Spec.Ports {
		ports = append(ports, port.Port)
	}

	return &ServiceInfo{
		Name:      service.Name,
		ClusterIP: service.Spec.ClusterIP,
		Ports:     ports,
		Labels:    service.Labels,
	}, nil
}

// Restart restarts service
func (ksm *KubernetesServiceManager) Restart(ctx context.Context, options ContainerCreateOptions) error {
	// Get existing service information
	_, err := ksm.Get(ctx, options.ServiceName)
	if err != nil {
		// If service does not exist, directly create new service
		if isNotFoundError(err) {
			_, createErr := ksm.Create(ctx, options.ServiceName, options.Port, options.Labels)
			if createErr != nil {
				return fmt.Errorf("failed to create service: %w", createErr)
			}
			return nil
		}
		return fmt.Errorf("failed to get service information: %w", err)
	}

	// Service exists, delete first then recreate
	if err := ksm.Delete(ctx, options.ServiceName); err != nil {
		return fmt.Errorf("failed to delete existing service: %w", err)
	}

	// Wait for service to be completely deleted
	if err := ksm.waitForServiceDeletion(ctx, options.ServiceName); err != nil {
		return fmt.Errorf("failed to wait for service deletion completion: %w", err)
	}

	// Recreate service (use original port and labels, if none then use passed labels)
	_, err = ksm.Create(ctx, options.ServiceName, options.Port, options.Labels)
	if err != nil {
		return fmt.Errorf("failed to recreate service: %w", err)
	}

	return nil
}

// waitForServiceDeletion waits for service to be completely deleted
func (ksm *KubernetesServiceManager) waitForServiceDeletion(ctx context.Context, serviceName string) error {
	const (
		maxRetries    = 15              // maximum retry count
		retryInterval = 1 * time.Second // retry interval
	)

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for service deletion was cancelled: %w", ctx.Err())
		default:
		}

		// Check if service still exists
		_, err := ksm.Get(ctx, serviceName)
		if err != nil {
			// If get fails and is NotFound error, deletion is successful
			if isNotFoundError(err) {
				return nil
			}
			// Other errors continue retrying
		} else {
			// Service still exists, continue waiting
		}

		time.Sleep(retryInterval)
	}
	return fmt.Errorf("waiting for service deletion timed out, exceeded %d seconds", maxRetries)
}

// KubernetesVolumeManager Kubernetes volume manager implementation
type KubernetesVolumeManager struct {
	Entry *k8s.Entry
}

func (kvm *KubernetesVolumeManager) List(ctx context.Context) ([]VolumeInfo, error) {
	return nil, fmt.Errorf("not implemented for Kubernetes")
}

func (kvm *KubernetesVolumeManager) Create(ctx context.Context, name string, driver string, labels map[string]string, options map[string]string) (VolumeInfo, error) {
	return VolumeInfo{}, fmt.Errorf("not implemented for Kubernetes")
}

func (kvm *KubernetesVolumeManager) Inspect(ctx context.Context, name string) (VolumeInfo, error) {
	return VolumeInfo{}, fmt.Errorf("not implemented for Kubernetes")
}

func (kvm *KubernetesVolumeManager) Remove(ctx context.Context, name string) error {
	return fmt.Errorf("not implemented for Kubernetes")
}

func (kvm *KubernetesVolumeManager) Prune(ctx context.Context) (int, error) {
	return 0, fmt.Errorf("not implemented for Kubernetes")
}

// isNotFoundError checks if it's a NotFound error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// Check if error message contains "not found" keyword
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "notfound")
}
