package taskflow

type Connection struct {
	id     uint64
	source *OutputSlot
	dest   *InputSlot
}

func NewConnection(id uint64, source *OutputSlot, dest *InputSlot) *Connection {
	return &Connection{
		id:     id,
		source: source,
		dest:   dest,
	}
}

func (conn *Connection) Apply() error {
	return conn.source.AddConnection(conn)
}

func (conn *Connection) Close() error {
	return conn.source.RemoveConnection(conn)
}

func (conn *Connection) Execute(message *Message) {
	conn.dest.Push(message)
}

func (conn *Connection) Push(message *Message) {
	conn.source.task.taskflow.GetScheduler().Push(conn.id, message)
}
