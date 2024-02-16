import { useEffect, useState } from "react"
import { useParams } from 'react-router-dom';

function Submit() {
    const { roomId } = useParams();

    const [parsedJson, setParsedJson] = useState(null)

    useEffect(()=>{
        const url = 'http://localhost:8000/event/'+roomId
        const eventSource = new EventSource(url)

        eventSource.onmessage = (event) => {
            console.log("onMessage")
            const newEvent = event.data
            try{
                setParsedJson(JSON.parse(newEvent))
            } catch(error) {
                console.log(`Error ${error}`)
            }
        }

        return () => {
            console.log("closing event")
            eventSource.close()
        }
    }, [])

    return(
        <div className="ScoreCard">
            {parsedJson && (
                <>
                    <h1>ScoreCard : {roomId}</h1>
                    <div className="Players">
                        {console.log(parsedJson)}
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