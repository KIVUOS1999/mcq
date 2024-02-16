import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';
import paths from "../constants/constants"

function JoinRoom() {
    const [playerID, setPlayerID] = useState('')
    const [roomID, setRoomID] = useState('')
    const navigate = useNavigate();

    const addPlayer = ()=> {
        const url = paths.base + paths.a_player + roomID +"/player/" +playerID
        fetch (url)
        .then(res => {
            if(!res.ok) {
                if(res.status === 400) {
                    res.json().then(errData => {
                        alert(`Error in JoinRoom: ${errData.response_code}|${errData.reason}`)
                    })
                    return
                }
                throw new Error(`HTTP: call error: ${res.status}`)
            }

            alert(`Player successfully added + ${roomID} : ${playerID}`)
            navigate(`/question/${roomID}/${playerID}`);
        })
        .catch(error => {
            console.error(error)
        })
    }

    const addParams = () => {
        const player = document.getElementById("player")
        const room = document.getElementById("room")

        setPlayerID(player.value)
        setRoomID(room.value)
    }

    useEffect(() => {
        if (roomID !== "" && playerID !== ""){
            addPlayer()
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