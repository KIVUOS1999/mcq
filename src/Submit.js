import { useEffect, useState } from "react"
import { useParams } from 'react-router-dom';
import constants from './constants/constants'

function Submit() {
    const { roomId } = useParams();

    const [parsedJson, setParsedJson] = useState(null)
    const [ended, setEnded] = useState(false)

    const fetchResult = ()=>{
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
            setEnded(true)
            return error
        })
    }

    useEffect(() => {
        if(ended) {
            clearInterval(interval)
            return
        }

        fetchResult()

        const interval = setInterval(() => {
            console.log("making fetch submit call")
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