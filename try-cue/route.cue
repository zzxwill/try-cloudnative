impo
cpuScaler: {
	condition: {
		if parameter["cpuPercent"] != _|_ {
			value: parameter.cpuPercent
		}
	}
}

parameter: {
	cpuPercent?: int
}
