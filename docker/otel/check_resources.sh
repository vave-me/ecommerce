#!/bin/bash

echo "=== Kubernetes Resource Analysis Report ==="
echo "Generated on: $(date)"
echo "==========================================="
echo

# Function to convert memory to Mi
convert_to_mi() {
    local value=$1
    if [[ $value =~ ^([0-9]+)Mi$ ]]; then
        echo "${BASH_REMATCH[1]}"
    elif [[ $value =~ ^([0-9]+)Gi$ ]]; then
        echo $((${BASH_REMATCH[1]} * 1024))
    elif [[ $value =~ ^([0-9]+)M$ ]]; then
        echo "${BASH_REMATCH[1]}"
    else
        echo "0"
    fi
}

# Function to convert CPU to millicores
convert_to_millicores() {
    local value=$1
    if [[ $value =~ ^([0-9]+)m$ ]]; then
        echo "${BASH_REMATCH[1]}"
    elif [[ $value =~ ^([0-9]+)$ ]]; then
        echo $((${BASH_REMATCH[1]} * 1000))
    else
        echo "0"
    fi
}

echo "1. SERVICES WITHOUT RESOURCE LIMITS OR REQUESTS:"
echo "================================================"
found_missing=false
find /home/szymon/classified/k8s -name "*.yaml" -type f | while read file; do
    if grep -q -E "kind:\s*(Deployment|StatefulSet|DaemonSet)" "$file" 2>/dev/null; then
        # Check if the deployment has containers
        if grep -q "containers:" "$file" 2>/dev/null; then
            # Count containers and resources blocks
            container_count=$(grep -c "- name:" "$file" | head -1)
            resource_count=$(grep -c "resources:" "$file" | head -1)
            
            if [ "$resource_count" -lt "$container_count" ] 2>/dev/null; then
                echo "⚠️  $file - Has containers without resource specifications"
                found_missing=true
            fi
        fi
    fi
done

if [ "$found_missing" = false ]; then
    echo "✅ All deployments have resource specifications"
fi

echo
echo "2. CONTAINERS WITH VERY HIGH RESOURCE LIMITS (potential cluster impact):"
echo "========================================================================"
echo "High CPU (>2 cores) or High Memory (>2Gi):"
echo

find /home/szymon/classified/k8s -name "*.yaml" -type f | while read file; do
    if grep -q -E "kind:\s*(Deployment|StatefulSet|DaemonSet)" "$file" 2>/dev/null; then
        # Extract container name and check limits
        if grep -A 10 "resources:" "$file" 2>/dev/null | grep -q "limits:"; then
            service_name=$(basename "$file" .yaml)
            
            # Check CPU limits
            cpu_limit=$(grep -A 5 "limits:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | head -1)
            if [ ! -z "$cpu_limit" ]; then
                cpu_millicores=$(convert_to_millicores "$cpu_limit")
                if [ "$cpu_millicores" -gt 2000 ] 2>/dev/null; then
                    echo "⚠️  HIGH CPU: $file - $service_name has CPU limit: $cpu_limit"
                fi
            fi
            
            # Check memory limits
            mem_limit=$(grep -A 5 "limits:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | head -1)
            if [ ! -z "$mem_limit" ]; then
                mem_mi=$(convert_to_mi "$mem_limit")
                if [ "$mem_mi" -gt 2048 ] 2>/dev/null; then
                    echo "⚠️  HIGH MEMORY: $file - $service_name has Memory limit: $mem_limit"
                fi
            fi
        fi
    fi
done

echo
echo "3. SERVICES WITH UNREASONABLE RESOURCE REQUESTS/LIMITS:"
echo "======================================================="
echo "Checking for potential issues..."
echo

find /home/szymon/classified/k8s -name "*.yaml" -type f | while read file; do
    if grep -q -E "kind:\s*(Deployment|StatefulSet|DaemonSet)" "$file" 2>/dev/null; then
        service_name=$(basename "$file" .yaml)
        
        # Extract requests and limits
        cpu_request=$(grep -A 5 "requests:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | head -1)
        cpu_limit=$(grep -A 5 "limits:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | head -1)
        mem_request=$(grep -A 5 "requests:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | head -1)
        mem_limit=$(grep -A 5 "limits:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | head -1)
        
        if [ ! -z "$cpu_request" ] && [ ! -z "$cpu_limit" ]; then
            req_millicores=$(convert_to_millicores "$cpu_request")
            lim_millicores=$(convert_to_millicores "$cpu_limit")
            
            # Check if limit is more than 10x the request
            if [ "$req_millicores" -gt 0 ] && [ "$lim_millicores" -gt 0 ] 2>/dev/null; then
                ratio=$((lim_millicores / req_millicores))
                if [ "$ratio" -gt 10 ] 2>/dev/null; then
                    echo "⚠️  LARGE CPU RATIO: $service_name - Request: $cpu_request, Limit: $cpu_limit (${ratio}x difference)"
                fi
            fi
        fi
        
        if [ ! -z "$mem_request" ] && [ ! -z "$mem_limit" ]; then
            req_mi=$(convert_to_mi "$mem_request")
            lim_mi=$(convert_to_mi "$mem_limit")
            
            # Check if limit is more than 4x the request
            if [ "$req_mi" -gt 0 ] && [ "$lim_mi" -gt 0 ] 2>/dev/null; then
                ratio=$((lim_mi / req_mi))
                if [ "$ratio" -gt 4 ] 2>/dev/null; then
                    echo "⚠️  LARGE MEMORY RATIO: $service_name - Request: $mem_request, Limit: $mem_limit (${ratio}x difference)"
                fi
            fi
        fi
    fi
done

echo
echo "4. RESOURCE CONFIGURATION SUMMARY BY SERVICE:"
echo "============================================"
echo
printf "%-30s %-15s %-15s %-15s %-15s\n" "SERVICE" "CPU REQUEST" "CPU LIMIT" "MEM REQUEST" "MEM LIMIT"
printf "%-30s %-15s %-15s %-15s %-15s\n" "-------" "-----------" "---------" "-----------" "---------"

# Services in 09-services
for file in /home/szymon/classified/k8s/09-services/*.yaml; do
    if [ -f "$file" ]; then
        service_name=$(basename "$file" .yaml)
        cpu_request=$(grep -A 5 "requests:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | tail -1)
        cpu_limit=$(grep -A 5 "limits:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | tail -1)
        mem_request=$(grep -A 5 "requests:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | tail -1)
        mem_limit=$(grep -A 5 "limits:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | tail -1)
        
        printf "%-30s %-15s %-15s %-15s %-15s\n" "$service_name" "${cpu_request:-N/A}" "${cpu_limit:-N/A}" "${mem_request:-N/A}" "${mem_limit:-N/A}"
    fi
done

echo
echo "OBSERVABILITY SERVICES:"
printf "%-30s %-15s %-15s %-15s %-15s\n" "-------" "-----------" "---------" "-----------" "---------"

# Observability services
for dir in 15-jaeger 16-grafana 17-otel 18-prometheus; do
    for file in /home/szymon/classified/k8s/$dir/*.yaml; do
        if [ -f "$file" ] && grep -q -E "kind:\s*(Deployment|StatefulSet)" "$file" 2>/dev/null; then
            service_name="$dir/$(basename "$file" .yaml)"
            cpu_request=$(grep -A 5 "requests:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | tail -1)
            cpu_limit=$(grep -A 5 "limits:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | tail -1)
            mem_request=$(grep -A 5 "requests:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | tail -1)
            mem_limit=$(grep -A 5 "limits:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | tail -1)
            
            printf "%-30s %-15s %-15s %-15s %-15s\n" "$service_name" "${cpu_request:-N/A}" "${cpu_limit:-N/A}" "${mem_request:-N/A}" "${mem_limit:-N/A}"
        fi
    done
done

echo
echo "5. TOTAL RESOURCE REQUESTS (Minimum cluster capacity needed):"
echo "==========================================================="

total_cpu_req=0
total_mem_req=0

find /home/szymon/classified/k8s -name "*.yaml" -type f | while read file; do
    if grep -q -E "kind:\s*(Deployment|StatefulSet|DaemonSet)" "$file" 2>/dev/null; then
        # Get replica count
        replicas=$(grep "replicas:" "$file" | awk '{print $2}' | head -1)
        replicas=${replicas:-1}
        
        # Get CPU request
        cpu_request=$(grep -A 5 "requests:" "$file" | grep "cpu:" | awk '{print $2}' | tr -d '"' | tail -1)
        if [ ! -z "$cpu_request" ]; then
            cpu_millicores=$(convert_to_millicores "$cpu_request")
            total_cpu_req=$((total_cpu_req + (cpu_millicores * replicas)))
        fi
        
        # Get memory request
        mem_request=$(grep -A 5 "requests:" "$file" | grep "memory:" | awk '{print $2}' | tr -d '"' | tail -1)
        if [ ! -z "$mem_request" ]; then
            mem_mi=$(convert_to_mi "$mem_request")
            total_mem_req=$((total_mem_req + (mem_mi * replicas)))
        fi
    fi
done | tail -1

echo "Total CPU Requests: $(echo "scale=2; $total_cpu_req / 1000" | bc) cores"
echo "Total Memory Requests: $(echo "scale=2; $total_mem_req / 1024" | bc) Gi"

echo
echo "=== END OF REPORT ==="