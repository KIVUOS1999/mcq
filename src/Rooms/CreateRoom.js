import { useEffect, useState } from "react"
import { useNavigate } from 'react-router-dom';
import paths from "../constants/constants"

function CreateRoom() {
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
        .then()
        .catch(error => {
            console.error(error)
        })
    }
    
    const createRoomAPI = ()=>{
        const url = paths.bigujjknouiknkjjase + paths.c_room
    
        fetch(url)
        .then(res => {
            if(!res.ok) 
            {
                if(res.status === 400) {
                    res.json().then(errData => {
                        alert(`Error in JoinRoom: ${errData.response_code}|${errData.reason}`)
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
            setPlayerID(player.value);
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
            addPlayer()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [roomID, playerID])
    
    return (
        <div className="CreateRoom">
            <h2>Create Room</h2> 
            Enter Name: <input id="player_name" type="text"></input><br/>
            <button onClick={createRoomAddPlayer}>Generate RoomCode</button><br/>
        </div>
    )
}

export default CreateRoom