# Security and Sandboxing

## Overview

AgentForge implements a comprehensive security model that protects against malicious components while enabling legitimate functionality. The system uses multi-layered security including component sandboxing, permission management, code signing, and runtime isolation.

## Security Architecture

### Multi-Layer Security Model
```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                         │
├─────────────────────────────────────────────────────────────┤
│  Application Layer                                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Component   │ │ Permission  │ │ Audit       │          │
│  │ Signing     │ │ Management  │ │ Logging     │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  Runtime Layer                                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Sandboxing  │ │ Resource    │ │ Network     │          │
│  │ & Isolation │ │ Limits      │ │ Filtering   │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
├─────────────────────────────────────────────────────────────┤
│  System Layer                                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ Container   │ │ Filesystem  │ │ Process     │          │
│  │ Security    │ │ Isolation   │ │ Isolation   │          │
│  └─────────────┘ └─────────────┘ └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

## Component Security

### Code Signing and Verification
```yaml
# Component signature metadata
signature:
  algorithm: "RSA-SHA256"
  public_key_id: "forge-official-2024"
  signature: "base64-encoded-signature"
  signed_at: "2024-01-15T10:30:00Z"
  signer: "forge-official@agentforge.dev"
  
verification:
  required: true
  trusted_signers:
    - "forge-official@agentforge.dev"
    - "company-team@company.com"
  allow_unsigned: false  # For development only
```

### Component Scanning
```go
type SecurityScanner struct {
    codeAnalyzer    *CodeAnalyzer
    dependencyCheck *DependencyChecker
    secretDetector  *SecretDetector
    malwareScanner  *MalwareScanner
}

func (ss *SecurityScanner) ScanComponent(ctx context.Context, component *Component) (*SecurityReport, error) {
    report := &SecurityReport{
        ComponentID: component.ID,
        ScanTime:    time.Now(),
        Findings:    []SecurityFinding{},
    }
    
    // Code analysis
    codeFindings, err := ss.codeAnalyzer.Analyze(ctx, component.SourceCode)
    if err != nil {
        return nil, fmt.Errorf("code analysis failed: %w", err)
    }
    report.Findings = append(report.Findings, codeFindings...)
    
    // Dependency vulnerability check
    depFindings, err := ss.dependencyCheck.CheckVulnerabilities(ctx, component.Dependencies)
    if err != nil {
        return nil, fmt.Errorf("dependency check failed: %w", err)
    }
    report.Findings = append(report.Findings, depFindings...)
    
    // Secret detection
    secretFindings, err := ss.secretDetector.Scan(ctx, component.SourceCode)
    if err != nil {
        return nil, fmt.Errorf("secret detection failed: %w", err)
    }
    report.Findings = append(report.Findings, secretFindings...)
    
    // Malware scanning
    malwareFindings, err := ss.malwareScanner.Scan(ctx, component.Binary)
    if err != nil {
        return nil, fmt.Errorf("malware scan failed: %w", err)
    }
    report.Findings = append(report.Findings, malwareFindings...)
    
    // Calculate overall risk score
    report.RiskScore = ss.calculateRiskScore(report.Findings)
    report.RiskLevel = ss.determineRiskLevel(report.RiskScore)
    
    return report, nil
}

type SecurityFinding struct {
    Type        FindingType
    Severity    Severity
    Title       string
    Description string
    Location    string
    Remediation string
    CVEID       string  // For vulnerability findings
    CWE         string  // Common Weakness Enumeration
}

type FindingType string
const (
    FindingTypeVulnerability FindingType = "vulnerability"
    FindingTypeSecret        FindingType = "secret"
    FindingTypeMalware       FindingType = "malware"
    FindingTypeCodeQuality   FindingType = "code_quality"
    FindingTypePermission    FindingType = "permission"
)
```

### Security Scanning Commands
```bash
# Scan component for security issues
forge security scan salesforce-lookup
forge security scan --all
forge security scan --composition enterprise-sales-stack

# Detailed security report
forge security report salesforce-lookup --format json
forge security report salesforce-lookup --include-remediation

# Vulnerability database updates
forge security update-db
forge security db-info

# Security policy compliance
forge security compliance salesforce-lookup --policy enterprise
forge security compliance --all --policy pci-dss
```

## Permission System

### Permission Model
```yaml
# Component permission specification
permissions:
  # Network permissions
  network:
    outbound:
      - domains: ["*.salesforce.com", "*.force.com"]
        ports: [443, 80]
        protocols: ["https", "http"]
      - domains: ["api.company.com"]
        ports: [443]
        protocols: ["https"]
        
    inbound:
      - ports: [8080]
        protocols: ["http"]
        source: "internal"
        
  # Filesystem permissions
  filesystem:
    read:
      - paths: ["/app/config", "/tmp"]
        recursive: true
      - paths: ["/etc/ssl/certs"]
        recursive: false
        
    write:
      - paths: ["/tmp", "/app/logs"]
        recursive: true
        
    execute:
      - paths: ["/app/bin"]
        recursive: false
        
  # Environment permissions
  environment:
    read:
      - variables: ["SALESFORCE_*", "API_*"]
      - variables: ["PATH", "HOME"]
        
    write: []  # No environment write permissions
    
  # System permissions
  system:
    processes:
      - spawn: false
      - signal: false
      
    resources:
      - cpu_limit: "1000m"
      - memory_limit: "512Mi"
      - disk_limit: "1Gi"
      
  # External service permissions
  services:
    - service: "salesforce_api"
      operations: ["read", "write"]
      rate_limit: "100/minute"
      
    - service: "email_service"
      operations: ["send"]
      rate_limit: "50/hour"
```

### Permission Management
```go
type PermissionManager struct {
    policies    map[string]*SecurityPolicy
    enforcer    *PermissionEnforcer
    auditor     *AuditLogger
}

func (pm *PermissionManager) CheckPermission(ctx context.Context, componentID string, permission Permission) error {
    component, err := pm.getComponent(ctx, componentID)
    if err != nil {
        return err
    }
    
    // Check if permission is granted
    if !pm.hasPermission(component, permission) {
        pm.auditor.LogPermissionDenied(ctx, componentID, permission)
        return &PermissionDeniedError{
            ComponentID: componentID,
            Permission:  permission,
            Message:     fmt.Sprintf("Permission %s denied for component %s", permission.Type, componentID),
        }
    }
    
    // Check rate limits
    if err := pm.checkRateLimit(ctx, componentID, permission); err != nil {
        return err
    }
    
    // Log permission usage
    pm.auditor.LogPermissionUsed(ctx, componentID, permission)
    
    return nil
}

func (pm *PermissionManager) GrantPermission(ctx context.Context, componentID string, permission Permission) error {
    // Validate permission request
    if err := pm.validatePermission(permission); err != nil {
        return fmt.Errorf("invalid permission: %w", err)
    }
    
    // Check if user has authority to grant permission
    if !pm.canGrantPermission(ctx, permission) {
        return &InsufficientPrivilegesError{
            Permission: permission,
            Message:    "Insufficient privileges to grant permission",
        }
    }
    
    // Grant permission
    component, err := pm.getComponent(ctx, componentID)
    if err != nil {
        return err
    }
    
    component.Permissions = append(component.Permissions, permission)
    
    if err := pm.updateComponent(ctx, component); err != nil {
        return fmt.Errorf("failed to update component permissions: %w", err)
    }
    
    pm.auditor.LogPermissionGranted(ctx, componentID, permission)
    
    return nil
}
```

### Permission Commands
```bash
# View component permissions
forge security permissions salesforce-lookup
forge security permissions salesforce-lookup --detailed

# Grant permissions
forge security grant salesforce-lookup --permission network.outbound
forge security grant salesforce-lookup --permission filesystem.read:/app/data
forge security grant salesforce-lookup --permission env.read:API_KEY

# Revoke permissions
forge security revoke salesforce-lookup --permission filesystem.write
forge security revoke salesforce-lookup --all-network

# Permission policies
forge security policy list
forge security policy apply enterprise --to salesforce-lookup
forge security policy create custom --from-file ./policy.yaml
```

## Sandboxing and Isolation

### Container-Based Sandboxing
```go
type Sandbox struct {
    containerID   string
    runtime       ContainerRuntime
    config        *SandboxConfig
    monitor       *ResourceMonitor
    networkFilter *NetworkFilter
}

type SandboxConfig struct {
    // Resource limits
    CPULimit      string  // "1000m"
    MemoryLimit   string  // "512Mi"
    DiskLimit     string  // "1Gi"
    NetworkLimit  string  // "10Mbps"
    
    // Filesystem isolation
    ReadOnlyPaths  []string
    ReadWritePaths []string
    TempDirs       []string
    
    // Network isolation
    AllowedDomains []string
    AllowedPorts   []int
    DNSServers     []string
    
    // Process isolation
    AllowProcessSpawn bool
    AllowSignals      bool
    MaxProcesses      int
    
    // Security options
    NoNewPrivileges bool
    DropCapabilities []string
    SeccompProfile   string
    AppArmorProfile  string
}

func (s *Sandbox) CreateSandbox(ctx context.Context, component *Component) error {
    // Create container configuration
    containerConfig := &container.Config{
        Image: component.RuntimeImage,
        Env:   s.buildEnvironment(component),
        Cmd:   component.Command,
        
        // Security settings
        User: "1000:1000", // Non-root user
        WorkingDir: "/app",
        
        // Resource limits
        Resources: container.Resources{
            CPUQuota:  s.parseCPULimit(s.config.CPULimit),
            Memory:    s.parseMemoryLimit(s.config.MemoryLimit),
            DiskQuota: s.parseDiskLimit(s.config.DiskLimit),
        },
        
        // Security options
        SecurityOpt: []string{
            "no-new-privileges:true",
            fmt.Sprintf("seccomp:%s", s.config.SeccompProfile),
            fmt.Sprintf("apparmor:%s", s.config.AppArmorProfile),
        },
    }
    
    // Create host configuration
    hostConfig := &container.HostConfig{
        // Network isolation
        NetworkMode: "custom",
        DNS:         s.config.DNSServers,
        
        // Filesystem isolation
        ReadonlyPaths: s.config.ReadOnlyPaths,
        Tmpfs: map[string]string{
            "/tmp": "rw,noexec,nosuid,size=100m",
        },
        
        // Capability dropping
        CapDrop: s.config.DropCapabilities,
        
        // Process limits
        PidsLimit: &s.config.MaxProcesses,
    }
    
    // Create and start container
    resp, err := s.runtime.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
    if err != nil {
        return fmt.Errorf("failed to create container: %w", err)
    }
    
    s.containerID = resp.ID
    
    if err := s.runtime.ContainerStart(ctx, s.containerID, types.ContainerStartOptions{}); err != nil {
        return fmt.Errorf("failed to start container: %w", err)
    }
    
    // Start monitoring
    go s.monitor.MonitorResources(ctx, s.containerID)
    go s.networkFilter.FilterTraffic(ctx, s.containerID)
    
    return nil
}
```

### Process Isolation
```go
type ProcessIsolator struct {
    namespaces []Namespace
    cgroups    *CgroupManager
    seccomp    *SeccompFilter
}

func (pi *ProcessIsolator) IsolateProcess(ctx context.Context, component *Component) (*IsolatedProcess, error) {
    // Create new namespaces
    namespaces := []Namespace{
        {Type: "pid"},   // Process ID namespace
        {Type: "net"},   // Network namespace
        {Type: "mnt"},   // Mount namespace
        {Type: "uts"},   // Hostname namespace
        {Type: "ipc"},   // Inter-process communication namespace
    }
    
    // Setup cgroups for resource control
    cgroupPath := fmt.Sprintf("/sys/fs/cgroup/agentforge/%s", component.ID)
    if err := pi.cgroups.CreateCgroup(cgroupPath); err != nil {
        return nil, fmt.Errorf("failed to create cgroup: %w", err)
    }
    
    // Configure resource limits
    limits := &ResourceLimits{
        CPUQuota:    component.Config.CPULimit,
        MemoryLimit: component.Config.MemoryLimit,
        DiskLimit:   component.Config.DiskLimit,
    }
    
    if err := pi.cgroups.SetLimits(cgroupPath, limits); err != nil {
        return nil, fmt.Errorf("failed to set resource limits: %w", err)
    }
    
    // Setup seccomp filter
    seccompFilter, err := pi.seccomp.CreateFilter(component.Config.AllowedSyscalls)
    if err != nil {
        return nil, fmt.Errorf("failed to create seccomp filter: %w", err)
    }
    
    // Start isolated process
    cmd := exec.CommandContext(ctx, component.Binary, component.Args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | 
                   syscall.CLONE_NEWMNT | syscall.CLONE_NEWUTS | 
                   syscall.CLONE_NEWIPC,
    }
    
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("failed to start isolated process: %w", err)
    }
    
    // Add process to cgroup
    if err := pi.cgroups.AddProcess(cgroupPath, cmd.Process.Pid); err != nil {
        cmd.Process.Kill()
        return nil, fmt.Errorf("failed to add process to cgroup: %w", err)
    }
    
    return &IsolatedProcess{
        Process:    cmd.Process,
        CgroupPath: cgroupPath,
        Namespaces: namespaces,
        Filter:     seccompFilter,
    }, nil
}
```

## Network Security

### Network Filtering
```go
type NetworkFilter struct {
    rules    []FilterRule
    monitor  *NetworkMonitor
    logger   *zap.Logger
}

type FilterRule struct {
    Direction   Direction  // inbound, outbound
    Protocol    string     // tcp, udp, icmp
    Source      string     // IP/CIDR or domain
    Destination string     // IP/CIDR or domain
    Port        int        // Port number
    Action      Action     // allow, deny, log
}

func (nf *NetworkFilter) FilterTraffic(ctx context.Context, containerID string) error {
    // Get container network namespace
    netNS, err := nf.getContainerNetNS(containerID)
    if err != nil {
        return fmt.Errorf("failed to get network namespace: %w", err)
    }
    
    // Apply iptables rules
    for _, rule := range nf.rules {
        iptablesRule := nf.convertToIptables(rule)
        if err := nf.applyIptablesRule(netNS, iptablesRule); err != nil {
            nf.logger.Error("Failed to apply iptables rule", 
                zap.String("rule", iptablesRule), 
                zap.Error(err))
        }
    }
    
    // Start traffic monitoring
    go nf.monitor.MonitorTraffic(ctx, netNS)
    
    return nil
}

func (nf *NetworkFilter) CheckConnection(ctx context.Context, componentID string, conn *Connection) error {
    // Check against filter rules
    for _, rule := range nf.rules {
        if nf.matchesRule(conn, rule) {
            switch rule.Action {
            case ActionAllow:
                nf.logger.Debug("Connection allowed", 
                    zap.String("component", componentID),
                    zap.String("destination", conn.Destination),
                    zap.Int("port", conn.Port))
                return nil
                
            case ActionDeny:
                nf.logger.Warn("Connection denied", 
                    zap.String("component", componentID),
                    zap.String("destination", conn.Destination),
                    zap.Int("port", conn.Port))
                return &ConnectionDeniedError{
                    ComponentID: componentID,
                    Connection:  conn,
                    Rule:        rule,
                }
                
            case ActionLog:
                nf.logger.Info("Connection logged", 
                    zap.String("component", componentID),
                    zap.String("destination", conn.Destination),
                    zap.Int("port", conn.Port))
                continue
            }
        }
    }
    
    // Default deny
    return &ConnectionDeniedError{
        ComponentID: componentID,
        Connection:  conn,
        Message:     "No matching allow rule found",
    }
}
```

### DNS Filtering
```go
type DNSFilter struct {
    allowedDomains []string
    blockedDomains []string
    resolver       *dns.Resolver
    cache          *DNSCache
}

func (df *DNSFilter) ResolveDomain(ctx context.Context, domain string) ([]net.IP, error) {
    // Check if domain is explicitly blocked
    if df.isBlocked(domain) {
        return nil, &DNSBlockedError{
            Domain: domain,
            Reason: "Domain is in blocklist",
        }
    }
    
    // Check if domain is allowed
    if !df.isAllowed(domain) {
        return nil, &DNSBlockedError{
            Domain: domain,
            Reason: "Domain not in allowlist",
        }
    }
    
    // Check cache first
    if ips := df.cache.Get(domain); ips != nil {
        return ips, nil
    }
    
    // Resolve domain
    ips, err := df.resolver.LookupIPAddr(ctx, domain)
    if err != nil {
        return nil, fmt.Errorf("DNS resolution failed: %w", err)
    }
    
    // Convert to net.IP slice
    result := make([]net.IP, len(ips))
    for i, ip := range ips {
        result[i] = ip.IP
    }
    
    // Cache result
    df.cache.Set(domain, result, 5*time.Minute)
    
    return result, nil
}
```

## Audit and Compliance

### Audit Logging
```go
type AuditLogger struct {
    logger    *zap.Logger
    storage   AuditStorage
    formatter AuditFormatter
}

type AuditEvent struct {
    Timestamp   time.Time
    EventType   string
    ComponentID string
    UserID      string
    Action      string
    Resource    string
    Result      string
    Details     map[string]interface{}
    IPAddress   string
    UserAgent   string
}

func (al *AuditLogger) LogSecurityEvent(ctx context.Context, event *AuditEvent) error {
    // Enrich event with context
    event.Timestamp = time.Now()
    event.UserID = al.getUserID(ctx)
    event.IPAddress = al.getClientIP(ctx)
    event.UserAgent = al.getUserAgent(ctx)
    
    // Format event
    formatted, err := al.formatter.Format(event)
    if err != nil {
        return fmt.Errorf("failed to format audit event: %w", err)
    }
    
    // Log to structured logger
    al.logger.Info("Security audit event",
        zap.String("event_type", event.EventType),
        zap.String("component_id", event.ComponentID),
        zap.String("action", event.Action),
        zap.String("result", event.Result),
        zap.Any("details", event.Details))
    
    // Store in audit storage
    if err := al.storage.Store(ctx, formatted); err != nil {
        al.logger.Error("Failed to store audit event", zap.Error(err))
        // Don't fail the operation due to audit storage failure
    }
    
    return nil
}

func (al *AuditLogger) LogPermissionDenied(ctx context.Context, componentID string, permission Permission) {
    event := &AuditEvent{
        EventType:   "permission_denied",
        ComponentID: componentID,
        Action:      "permission_check",
        Resource:    permission.Resource,
        Result:      "denied",
        Details: map[string]interface{}{
            "permission_type": permission.Type,
            "permission_value": permission.Value,
        },
    }
    
    al.LogSecurityEvent(ctx, event)
}
```

### Compliance Reporting
```bash
# Generate compliance reports
forge security compliance-report --standard pci-dss
forge security compliance-report --standard soc2 --period 2024-01
forge security compliance-report --custom ./compliance-rules.yaml

# Audit log analysis
forge security audit-logs --component salesforce-lookup --since 7d
forge security audit-logs --event-type permission_denied --format csv
forge security audit-logs --user john.doe --actions grant,revoke

# Security metrics
forge security metrics --dashboard
forge security metrics --component salesforce-lookup --period 30d
```

## Security Policies

### Policy Definition
```yaml
# security-policy.yaml
apiVersion: "forge.dev/v1"
kind: "SecurityPolicy"
metadata:
  name: "enterprise-policy"
  version: "1.0.0"
  description: "Enterprise security policy for production environments"

spec:
  # Component requirements
  component_requirements:
    signing:
      required: true
      trusted_signers:
        - "forge-official@agentforge.dev"
        - "company-security@company.com"
      
    scanning:
      required: true
      max_risk_score: 7.0
      block_on_high_severity: true
      
    verification:
      check_dependencies: true
      check_licenses: true
      allowed_licenses: ["MIT", "Apache-2.0", "BSD-3-Clause"]
  
  # Runtime restrictions
  runtime_restrictions:
    sandboxing:
      required: true
      container_runtime: "gvisor"  # or "runc", "kata"
      
    resource_limits:
      max_cpu: "2000m"
      max_memory: "1Gi"
      max_disk: "5Gi"
      max_network: "100Mbps"
      
    network_policy:
      default_action: "deny"
      allowed_domains:
        - "*.company.com"
        - "api.salesforce.com"
        - "api.openai.com"
      blocked_domains:
        - "*.malicious.com"
        - "*.suspicious.net"
      
    filesystem_policy:
      read_only_paths:
        - "/etc"
        - "/usr"
        - "/bin"
      writable_paths:
        - "/tmp"
        - "/app/logs"
        - "/app/cache"
  
  # Audit requirements
  audit_requirements:
    log_all_actions: true
    log_network_connections: true
    log_file_access: true
    retention_period: "2y"
    
  # Compliance mappings
  compliance:
    pci_dss:
      - requirement: "2.2.4"
        description: "Configure system security parameters"
        controls: ["sandboxing", "resource_limits"]
      - requirement: "10.2.1"
        description: "Log all individual user accesses"
        controls: ["audit_requirements"]
        
    soc2:
      - control: "CC6.1"
        description: "Logical and physical access controls"
        controls: ["network_policy", "filesystem_policy"]
```

### Policy Enforcement
```go
type PolicyEnforcer struct {
    policies map[string]*SecurityPolicy
    scanner  *SecurityScanner
    auditor  *AuditLogger
}

func (pe *PolicyEnforcer) EnforcePolicy(ctx context.Context, componentID string, policyName string) error {
    policy, exists := pe.policies[policyName]
    if !exists {
        return fmt.Errorf("policy %s not found", policyName)
    }
    
    component, err := pe.getComponent(ctx, componentID)
    if err != nil {
        return err
    }
    
    // Check component requirements
    if err := pe.checkComponentRequirements(ctx, component, policy.ComponentRequirements); err != nil {
        pe.auditor.LogSecurityEvent(ctx, &AuditEvent{
            EventType:   "policy_violation",
            ComponentID: componentID,
            Action:      "policy_check",
            Result:      "failed",
            Details: map[string]interface{}{
                "policy": policyName,
                "violation": err.Error(),
            },
        })
        return fmt.Errorf("component requirements not met: %w", err)
    }
    
    // Apply runtime restrictions
    if err := pe.applyRuntimeRestrictions(ctx, component, policy.RuntimeRestrictions); err != nil {
        return fmt.Errorf("failed to apply runtime restrictions: %w", err)
    }
    
    // Setup audit requirements
    if err := pe.setupAuditRequirements(ctx, component, policy.AuditRequirements); err != nil {
        return fmt.Errorf("failed to setup audit requirements: %w", err)
    }
    
    pe.auditor.LogSecurityEvent(ctx, &AuditEvent{
        EventType:   "policy_applied",
        ComponentID: componentID,
        Action:      "policy_enforcement",
        Result:      "success",
        Details: map[string]interface{}{
            "policy": policyName,
        },
    })
    
    return nil
}
```

## Security Commands

### Security Management CLI
```bash
# Security scanning
forge security scan <component>
forge security scan --all --policy enterprise
forge security scan --composition sales-stack

# Policy management
forge security policy list
forge security policy show enterprise-policy
forge security policy apply enterprise-policy --to salesforce-lookup
forge security policy validate ./custom-policy.yaml

# Permission management
forge security permissions <component>
forge security grant <component> --permission <permission>
forge security revoke <component> --permission <permission>

# Audit and compliance
forge security audit --component <component> --since 7d
forge security compliance-check --policy pci-dss
forge security report --format pdf --output security-report.pdf

# Trust management
forge security trust add github.com/company/tools
forge security trust list
forge security trust revoke github.com/untrusted/repo

# Incident response
forge security incident create --component <component> --severity high
forge security incident list --status open
forge security incident close <incident-id>
```

## Best Practices

### Security Guidelines
1. **Principle of Least Privilege**: Grant minimal permissions required
2. **Defense in Depth**: Multiple security layers
3. **Zero Trust**: Verify everything, trust nothing
4. **Regular Scanning**: Continuous security monitoring
5. **Audit Everything**: Comprehensive logging and monitoring

### Development Security
1. **Secure by Default**: Components start with minimal permissions
2. **Security Testing**: Include security tests in CI/CD
3. **Dependency Management**: Regular vulnerability scanning
4. **Code Review**: Security-focused code reviews
5. **Threat Modeling**: Consider security implications in design

### Operational Security
1. **Regular Updates**: Keep security policies current
2. **Incident Response**: Prepared response procedures
3. **Monitoring**: Real-time security monitoring
4. **Training**: Security awareness for developers
5. **Compliance**: Regular compliance audits