import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';
import paths from "../constants/constants"

function CreateRoom() {
    const [playerID, setPlayerID] = useState('')
    const [roomID, setRoomID] = useState('')
    const [time, setTime] = useState('')

    const navigate = useNavigate();

    const navigateToWaitingRoom = ()=> {
        navigate(`/lobby/${roomID}/${playerID}/1/${time}`);
    }
    
    const createRoomAPI = ()=>{
        const url = paths.base + paths.c_room
    
        fetch(url)
        .then(res => {
            if(!res.ok) 
            {
                if(res.status === 400) {
                    res.json().then(errData => {
                        alert(`Error in CreateRoom: ${errData.response_code}|${errData.reason}`)
                    })
                    return
                }
            throw new Error(`HTTP: call error: ${res.status}`)
            } 
    
            return res.json()
        })
        .then(data => {
            const roomNo = data.room_id
            setRoomID(roomNo)

            const player = document.getElementById("player_name")
            const player_time = document.getElementById("time")

            setPlayerID(player.value);
            setTime(player_time.value)
        })
        .catch(error => {
            console.error(error)
        })
    }
    
    const createRoomAddPlayer = ()=> {
        createRoomAPI()
    }

    useEffect(() => {
        if (roomID !== "" && playerID !== ""){
            navigateToWaitingRoom()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [roomID, playerID])
    
    return (
        <div className="CreateRoom">
            <h2>Create Room</h2> 
            Enter Name: <input id="player_name" type="text"></input><br/>
            Enter mins: <input id="time" type="text"></input><br/>
            <button onClick={createRoomAddPlayer}>Generate RoomCode</button><br/>
        </div>
    )
}

export default CreateRoom