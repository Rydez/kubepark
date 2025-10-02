.PHONY: help setup build deploy-park deploy-carousel deploy-restroom status logs clean open-grafana restart-monitoring reload-dashboards reset-dashboards

# Default target
help:
	@echo "🎢 KubePark Game Commands"
	@echo ""
	@echo "Setup & Build:"
	@echo "  setup           Create Kind cluster and start monitoring stack"
	@echo "  build           Build and push the game image"
	@echo ""
	@echo "Game Commands:"
	@echo "  deploy-park     Start the park (begins the game!)"
	@echo "  deploy-carousel Deploy carousel attraction"
	@echo "  deploy-restroom Deploy restroom attraction"
	@echo ""
	@echo "Monitoring:"
	@echo "  status          Show current park status"
	@echo "  logs            View park logs"
	@echo "  open-grafana    Open Grafana dashboard"
	@echo "  restart-monitoring Restart monitoring stack"
	@echo "  restart-alloy   Restart Alloy log collection"
	@echo "  upgrade-alloy   Upgrade Alloy with latest configuration"
	@echo "  reload-dashboards Reload Grafana dashboards"
	@echo "  reset-dashboards Complete dashboard reset (clears cache)"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean           Clean up everything"

# Setup infrastructure
setup:
	@echo "🏗️  Setting up KubePark infrastructure..."
	@echo "Creating Kind cluster with local registry..."
	@cd cluster && ./kind-with-registry.sh
	@echo "Starting monitoring stack..."
	@docker-compose up -d
	@echo "Waiting for monitoring stack to be ready..."
	@sleep 10
	@echo "Setting up Kubernetes resources..."
	@kubectl apply -f k8s/base/
	@echo "Adding Grafana Helm repository..."
	@helm repo add grafana https://grafana.github.io/helm-charts
	@helm repo update
	@echo "Deploying log collection (Alloy) using Helm..."
	@helm upgrade --install alloy grafana/alloy --namespace alloy --values k8s/alloy-values.yaml --wait
	@echo "✅ Setup complete! Run 'make build' next."

# Build and push image
build:
	@echo "🔨 Building KubePark image..."
	@cd .. && docker build -t localhost:5001/kubepark .
	@echo "📤 Pushing image to local registry..."
	@docker push localhost:5001/kubepark:latest
	@echo "✅ Build complete! Run 'make deploy-park' to start the game."

# Deploy park service
deploy-park:
	@echo "🎪 Starting the park..."
	@kubectl apply -f k8s/park.yaml
	@echo "⏳ Waiting for park to be ready..."
	@kubectl wait --for=condition=available --timeout=60s deployment/park -n park
	@echo "✅ Park is open! Check status with 'make status'"

# Deploy carousel attraction
deploy-carousel:
	@echo "🎠 Deploying carousel attraction..."
	@kubectl apply -f k8s/attractions/carousel.yaml
	@echo "⏳ Waiting for carousel to be ready..."
	@kubectl wait --for=condition=available --timeout=60s deployment/carousel -n attractions
	@echo "✅ Carousel is operational!"

# Deploy restroom attraction
deploy-restroom:
	@echo "🚻 Deploying restroom attraction..."
	@kubectl apply -f k8s/attractions/restroom.yaml
	@echo "⏳ Waiting for restroom to be ready..."
	@kubectl wait --for=condition=available --timeout=60s deployment/restroom -n attractions
	@echo "✅ Restroom is operational!"


# Show status
status:
	@echo "🎢 KubePark Status"
	@echo "=================="
	@echo ""
	@echo "📊 Cluster Status:"
	@kubectl get nodes
	@echo ""
	@echo "🎪 Park Services:"
	@kubectl get pods,svc -n park
	@echo ""
	@echo "🎠 Attractions:"
	@kubectl get pods,svc -n attractions
	@echo ""
	@echo "👥 Guests:"
	@kubectl get jobs,pods -n guests 2>/dev/null || echo "No guest jobs currently running (managed by park service)"
	@echo ""
	@echo "📈 Monitoring:"
	@docker-compose ps
	@echo ""
	@echo "📋 Log Collection:"
	@kubectl get pods -n alloy -l app.kubernetes.io/name=alloy

# View logs
logs:
	@echo "📋 Recent park logs:"
	@kubectl logs -n park deployment/park --tail=20 || echo "Park not deployed yet"

# Open Grafana
open-grafana:
	@echo "📊 Opening Grafana dashboard..."
	@echo "URL: http://localhost:3000"
	@echo "Username: admin"
	@echo "Password: admin"
	@open http://localhost:3000 2>/dev/null || echo "Please open http://localhost:3000 manually"

# Restart monitoring stack
restart-monitoring:
	@echo "🔄 Restarting monitoring stack..."
	@echo "Stopping monitoring services..."
	@docker-compose down
	@echo "Starting monitoring services..."
	@docker-compose up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 15
	@echo "✅ Monitoring stack restarted!"
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Prometheus: http://localhost:9090"
	@echo "Loki: http://localhost:3100"

# Force complete dashboard refresh (nuclear option)
reset-dashboards:
	@echo "🔄 Completely resetting Grafana dashboards..."
	@echo "Stopping Grafana..."
	@docker-compose stop grafana
	@echo "Removing Grafana container to clear cache..."
	@docker-compose rm -f grafana
	@echo "Removing Grafana data volume to clear all cached dashboards..."
	@docker volume rm walkthrough_grafana_data
	@echo "Starting fresh Grafana with clean state..."
	@docker-compose up -d grafana
	@echo "⏳ Waiting for Grafana to initialize completely..."
	@sleep 25
	@echo "✅ Grafana reset complete! All dashboards will be freshly loaded."
	@echo "Login at http://localhost:3000 with admin/admin"

# Clean up everything
clean:
	@echo "🧹 Cleaning up KubePark..."
	@kubectl delete namespace park attractions guests alloy --ignore-not-found=true
	@kubectl delete clusterrolebinding park-cluster-admin guests-cluster-admin --ignore-not-found=true
	@helm uninstall alloy --namespace alloy --ignore-not-found 2>/dev/null || true
	@kubectl delete namespace alloy --ignore-not-found=true
	@docker-compose down -v
	@kind delete cluster --name kubepark-cluster
	@docker rm -f kind-registry 2>/dev/null || true
	@echo "✅ Cleanup complete!"

# Advanced commands for power users
restart-park:
	@kubectl rollout restart deployment/park -n park

restart-attraction:
	@kubectl rollout restart deployment/carousel -n attractions || true
	@kubectl rollout restart deployment/restroom -n attractions || true

# Restart Alloy log collection
restart-alloy:
	@echo "🔄 Restarting Alloy log collection..."
	@kubectl rollout restart daemonset/alloy -n alloy
	@echo "✅ Alloy restarted!"

# Upgrade Alloy with latest configuration
upgrade-alloy:
	@echo "⬆️  Upgrading Alloy configuration..."
	@helm upgrade alloy grafana/alloy --namespace alloy --values k8s/alloy-values.yaml --wait
	@echo "✅ Alloy upgraded!"


