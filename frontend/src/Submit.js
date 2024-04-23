import { useEffect, useState } from "react"
import { useParams } from 'react-router-dom';
import {connect} from "nats.ws";
import constants from './constants/constants'

function Submit() {
    const { roomId } = useParams();

    const [parsedJson, setParsedJson] = useState(null)
    const [messages, setMessages] = useState("");
    const [natsSubs, setNatsSubs] = useState()

    useEffect(() => {
        const connectToNats = async () => {
            const nc = await connect({ servers: "ws://localhost:8222" });
            const sub = nc.subscribe(roomId);
            
            setNatsSubs(sub)

            for await (const msg of sub) {
                setMessages(msg);
            }
        }

        connectToNats();
    }, [roomId]);

    // nats player addition and starting game
    useEffect(()=>{
        const data = messages._rdata
        const natsData = new TextDecoder().decode(data)

        console.log(natsData)

        if (natsData && natsData != ""){ 
            setParsedJson(JSON.parse(natsData))
        }
    }, [messages])

    // api call to ask server to pusblish result
    useEffect(() => {
        const url = constants.base + "/get_result/" + roomId
        fetch (url)
        .then(res => {
            if(!res.ok) {
                if(!res.ok) {
                    res.json().then(errData => {
                        alert(`Error in Sending the message`)
                    })
                    return
                }
                throw new Error(`HTTP: call error: ${res.status}`)
            }
        })
        .catch(error => {
            console.error(error)
        })
    }, [natsSubs])

    return(
        <div className="ScoreCard">
            {parsedJson && (
                <>
                    <h1>ScoreCard : {roomId}</h1>
                    <div className="Players">
                        <ul>
                            {
                                parsedJson.score_card.map((item, index) => {
                                    return <li key={index}>{item.player_id} : {
                                        item.correct_answer === -1 ? "Yet to answer" : item.correct_answer
                                    }</li>
                                })
                            }
                        </ul>
                    </div>
                </>
            )}
        </div>
    )
}

export default Submit