package main

func debug(conn *ExtendedConnection) {
	conn.Connection.Emit("debug", "")
}
