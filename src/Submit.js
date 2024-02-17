import { useEffect, useState } from "react"
import { useParams } from 'react-router-dom';
import constants from './constants/constants'

function Submit() {
    const { roomId } = useParams();

    const [parsedJson, setParsedJson] = useState(null)

    const fetchResult = ()=>{
        {console.log("asking for scorebord")}
        const url = `${constants.base}/event/${roomId}`
        fetch(url)
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP fetch error: ${res.status}`)
            }

            return res.json()
        }).then(jsonData => {
            setParsedJson(jsonData)
        }).catch(error => {
            console.log(error)
        })
    }

    useEffect(() => {
        fetchResult()

        const interval = setInterval(() => {
            fetchResult()
        }, 5000)

        return() => clearInterval(interval)
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