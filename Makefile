create_namespaces:
	@temporal operator namespace create --namespace scenario1
	@temporal operator namespace create --namespace scenario6
delete_namespace:
	@tempolar operator namespace delete --namespace scenario1
	@tempolar operator namespace delete --namespace scenario2
	@tempolar operator namespace delete --namespace scenario6

delete_workflow:
	@tdbg --namespace scenario1 workflow deelte --workflow-id 6738999637