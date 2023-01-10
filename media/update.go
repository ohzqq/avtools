package media

type UpdateCmd struct {
	*Media
}

func (up UpdateCmd) Run() error {
	if up.MetaChanged {
		file := up.Input.NewName()
		file.Tmp(up.DumpIni())
		file.Run()
		tmp := file.file.Name()
		cmd := up.Command()
		cmd.Input.FFMeta(tmp)

		cmd.Output.Set("c", "copy")
		name := up.Input.NewName().Prefix("updated-").Join()
		cmd.Output.Ext(up.Input.Ext).Name(name).Pad("")

		c := cmd.Compile()
		//fmt.Println(c.String())

		err := c.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
