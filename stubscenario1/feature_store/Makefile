build:
	@docker build -f feature_store.Dockerfile --tag=feature_store:latest .
run:
	@docker run -d --rm --name feature_store -p 8000:8000 feature_store:latest