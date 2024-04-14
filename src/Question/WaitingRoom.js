import {connect} from "nats.ws";
import {useEffect,useState,useRef} from "react";
import { useParams, useNavigate } from 'react-router-dom';
import paths from "../constants/constants"

const WaitingRoom = (props)  => {
    const { roomId, playerId, isAdmin } = useParams();

    const [messages, setMessages] = useState([]);
    const [natsSubs, setNatsSubs] = useState()
    const [playersInLobby, setPlayersInLobby] = useState("")

    const isAddPlayerApiCalled = useRef(false)

    const startGame = "startGame:"

    const navigate = useNavigate()

    // settingup nats consumer
    useEffect(() => {
        const connectToNats = async () => {
            const nc = await connect({ servers: "ws://localhost:4222" });
            const sub = nc.subscribe(roomId);
            
            setNatsSubs(sub)

            for await (const msg of sub) {
                setMessages(msg);
            }
        }

        connectToNats();
    }, [props.roomId]);

    // nats player addition and starting game
    useEffect(()=>{
        const data = messages._rdata
        const players = new TextDecoder().decode(data)

        if (players.startsWith(startGame)) {
            console.log("game has started")

            navigate(`/question/${roomId}/${playerId}/${isAdmin}`)
            return
        }

        setPlayersInLobby(players)
    }, [messages])

    useEffect(()=>{
        if(natsSubs !== undefined && !isAddPlayerApiCalled.current){
            console.log("Calling add player API")
            addPlayer()
            isAddPlayerApiCalled.current = true
        }
    }, [natsSubs])

    const addPlayer = ()=> {
        const url = paths.base + paths.a_player + roomId +"/player/" +playerId + "/admin/1"
        fetch (url)
        .then(res => {
            if(!res.ok) {
                if(res.status === 400) {
                    res.json().then(errData => {
                        alert(`Error in AddPlayer: ${errData.response_code}|${errData.reason}`)
                    })
                    return
                }
                throw new Error(`HTTP: call error: ${res.status}`)
            }
        })
        .catch(error => {
            console.error(error)
        })
    }

    // admin start button action
    const sendStartMessage = () => {
        const url = paths.base + paths.s_game + roomId + "/endtime/" + paths.endTime
        fetch(url)
        .then(res => {
            if(!res.ok) {
                if(res.status === 400) {
                    console.error("Error in nats start game")
                    return
                } 
            }
        })
        .catch(error => {
            console.error(error)
        })
    } 

    return (
        <div className="Waiting_room">
            <h1>{`Waiting room ${roomId}: ${playerId}`}</h1>
            <p>{`Joined Players in room: ${playersInLobby}`}</p>
            {isAdmin == 1? <button onClick={sendStartMessage}>Start</button>:""}
        </div>
    )
}

export default WaitingRoom