import { useEffect, useRef, useState } from "react"
import { useParams } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import {connect} from "nats.ws";

import paths from '../constants/constants'
import QuestionOptions from './QuestionOptions'

function Question() {
    const { roomId, playerId, time } = useParams();
    const navigate = useNavigate()
    const [questionSet, setQuestionSet] = useState({
        Question: "",
        Options: [],
        Id:""
    })
    
    console.log("question>>", time)

    const [lenQues, setLenQues] = useState(0) 
    const [question, setQuestion] = useState()
    const [idx, setIdx] = useState(0)
    const [selectedOption, setSelectedOption] = useState(0);
    const [answer, setAnswer] = useState({})
    const [completed, setCompleted] = useState(false)
    const [countdownTimer, setContDownTimer] = useState(time)

    const ns = useRef(false)

    const [messages, setMessages] = useState();

    // nats server
    const connectToNats = async () => {
        console.log("connecting to nats server from question", roomId)
        
        const nc = await connect({ servers: "ws://localhost:8222" });
        const sub = nc.subscribe(roomId);

        for await (const msg of sub) {
            setMessages(msg);
        }

        ns.current = true
    }

    // getting questions set from server
    const getQuestionSet = () => {
        const url = `${paths.base}${paths.g_mcq}`
        fetch(url)
        .then(
            res => {
                if(!res.ok) {
                    throw new Error(`HTTP fetch error Status: ${res.status}`)
                }

                return res.json()
            }
        ).then(data => {
            setQuestion(data)
            setLenQues((data.question_set).length)
        }).catch(error => {
            console.error(`Error fetching data, ${error}`)
        })
    }

    // submitting answers for the player
    const submitAnswersCall = (jsonBody) => {
        const url = `${paths.base}${paths.a_answer}`
        fetch(url, {
            method: "POST",
            body: jsonBody
        })
        .then(
            res => {
                if(res.status !== 201) {
                    if(res.status === 400) {
                        res.json().then(errData => {
                            alert(`Error in Submitting answers: ${errData.response_code}|${errData.reason}`)
                        })
                        return
                    }
                    throw new Error(`HTTP fetch error Status: ${res.status}`)
                }
                else{
                    navigate(`/submit/${roomId}`);
                }
            }
        )
        .catch(error => {
            console.error(`Error fetching data, ${error}`)
        })
    }

    // Use Effects
    useEffect(() => {
        if (messages && new TextDecoder().decode(messages._rdata) === "submit_game") {
            submit()
        }
    }, [messages])
    
    useEffect(()=>{
        getQuestionSet()
        if (ns.current == false){
            connectToNats()
        }

        var x = setInterval(()=>{
            setContDownTimer(prevTime => prevTime-1)
        }, 1000)

        return () => clearInterval(x);
        // eslint-disable-next-line
    }, [])

    useEffect(() => {
        if (question && question.question_set) {
            setQuestionSet({
                Question: question.question_set[idx].question,
                Options: question.question_set[idx].options,
                Id: question.question_set[idx].id
            })
        }
    }, [idx, question])

    useEffect(() => {
        if (Object.keys(answer).length === 0) {
            return
        }
        submit()
    }, [completed])

    // Button Functions
    let nxt = (e) => {
        if (idx+1 === lenQues) {
            setCompleted(true)
            return
        }
        
        setIdx(idx+1)

        const questionIndex = e.target.id.toString()
        const so = selectedOption+1

        setAnswer({...answer, [questionIndex]: so})
    }

    let submit = () => {
        let preparedJson = {
            "room_id": roomId,
            "player_id": playerId,
            "answer_set": answer
        }

        const jsonPayload = JSON.stringify(preparedJson)

        console.log(jsonPayload)
        submitAnswersCall(jsonPayload)
    }

    return (
        <div className="Question Block">
            <h2>{`Welcome ${playerId} in Room: ${roomId}`}</h2>
            {questionSet.Question !== "" && (
                <>
                    <div className="Question">{questionSet.Id}. {questionSet.Question}</div>
                    <ul>
                        <QuestionOptions qs={questionSet} si={selectedOption} ssi={setSelectedOption}/>
                    </ul>
                    <div className="timer">{countdownTimer}</div>
                    {completed ? <button onClick={submit}>Submit</button>  : <button onClick={nxt} id={questionSet.Id}>NextQuestion</button>}
                </>
            )}
        </div>
    )
}

export default Question