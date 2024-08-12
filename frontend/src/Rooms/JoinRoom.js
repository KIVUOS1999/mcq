import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';

function JoinRoom() {
    const [playerID, setPlayerID] = useState('')
    const [roomID, setRoomID] = useState('')
    
    const navigate = useNavigate();

    const navigateToWaitingRoom = ()=> {
        navigate(`/lobby/${roomID}/${playerID}/0/nan`);
    }

    const addParams = () => {
        const player = document.getElementById("player")
        const room = document.getElementById("room")

        setPlayerID(player.value)
        setRoomID(room.value)
    }

    useEffect(() => {
        if (roomID !== "" && playerID !== ""){
            navigateToWaitingRoom()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [roomID, playerID])

    return(
        <div className="JoinRoom">
            <input className="Text Name" type="text" id="player" placeholder="Name" required></input>
            <input className="Text RoomCode" type="text" id="room" placeholder="Room Code" required></input>
            <button className="Btn ActionBtn" onClick={addParams}>Join Room</button>
        </div>
    )
}

export default JoinRoom