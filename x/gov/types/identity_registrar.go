package types

func WrapInfos(infos map[string]string) []IdentityInfoEntry {
	wrappedInfos := []IdentityInfoEntry{}
	for key, value := range infos {
		wrappedInfos = append(wrappedInfos, IdentityInfoEntry{
			Key:  key,
			Info: value,
		})
	}
	return wrappedInfos
}

func UnwrapInfos(wrappedInfos []IdentityInfoEntry) map[string]string {
	infos := make(map[string]string)
	for _, wi := range wrappedInfos {
		infos[wi.Key] = wi.Info
	}
	return infos
}
