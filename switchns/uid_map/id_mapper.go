package uid_map

type IDMapper interface {
	MapContainerName(container_name string) (string, error)
}
