# Bugs you can meet

## clientCtx.GetFromAddress() is returning empty address, and more in general can't read command flags

This can be related to you didn't run ReadTxCommandFlags before the command.
```go
    clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
    if err != nil {
        return err
    }
```