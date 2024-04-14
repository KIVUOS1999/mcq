import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';

function JoinRoom() {
    const [playerID, setPlayerID] = useState('')
    const [roomID, setRoomID] = useState('')
    
    const navigate = useNavigate();

    const navigateToWaitingRoom = ()=> {
        navigate(`/lobby/${roomID}/${playerID}/0`);
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
            <h2>Join Room</h2>
            Enter Name: <input type="text" id="player"></input><br/>
            Room Code: <input type="text" id="room"></input><br/>
            <button onClick={addParams}>Join Room</button>
        </div>
    )
}

export default JoinRoom