# Code Standard

## Folder Structure

```
├── bin/ # bin can handle for different tasks (such gomock)
│
├── cmd/ # handle all the dependencies and the entry point for different apps
│   ├── server.go
│   └── cli.go
│
├── internal/
│   ├── adapter/ # to adapte cli or http input
│   ├── repo/    # data access, it could be database, http or memory store
│   └── app/     # Business rules
│
├── pkg/
│
```

## Do

- Add dependency in consumer side, name it as `dependency.go`
- Use uber gomock to generate mock `dependency_mock.go` in the same place with `dependency.go`
- Mapping test if possible
- Put the struct in `struct.go` file

## Do not

- Use the name like `NoteDTO`. The struct name should be represent the idea of layer or domain. e.g. `HttpNoteResponse`, `DBNote`