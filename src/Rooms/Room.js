import { useState } from "react"
import CreateRoom from "./CreateRoom"
import JoinRoom from "./JoinRoom"

function Room() {
    const [createRoom, setCreateRoom] = useState(false)
    return(
        <div className="RoomBody">
            <h1>MCQ</h1>

            <button onClick={()=>{
                setCreateRoom(true)
            }} disabled={createRoom}>Create Room</button>
            
            <button onClick={()=>{
                setCreateRoom(false)
            }} disabled={!createRoom}>Join Room</button>

            <hr/>
            
            {createRoom ? <CreateRoom/> : <JoinRoom/>}
            
        </div>
    )
}

export default Room