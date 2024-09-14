import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';
import { toast, Toaster } from "sonner";

function isValidInput(value) {
    const isOK = /^[a-zA-Z0-9]+$/.test(value); // Regex for letters and numbers only
    console.log(isOK)

    return isOK
}

function JoinRoom() {
    const [playerID, setPlayerID] = useState('')
    const [roomID, setRoomID] = useState('')
    
    const navigate = useNavigate();

    const navigateToWaitingRoom = ()=> {
        navigate(`/lobby/${roomID}/${playerID}/0/nan`, {replace: true});
    }

    const addParams = () => {
        const player = document.getElementById("player")
        const room = document.getElementById("room")

        if (!isValidInput(player.value)){
            toast.error("player name not valid should be number and alphabets")
            return
        }

        // if (!isValidInput(room.value)) {
        //     toast.error("room id not valid should be number and alphabets")
        //     return
        // }
        
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
            <Toaster/>
        </div>
    )
}

export default JoinRoom