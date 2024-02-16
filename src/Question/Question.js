import { useEffect, useState } from "react"
import { useParams } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';

import paths from '../constants/constants'
import QuestionOptions from './QuestionOptions'

function Question() {
    const { roomId, playerId } = useParams();
    const navigate = useNavigate()
    const [questionSet, setQuestionSet] = useState({
        Question: "",
        Options: [],
        Id:""
    })
    
    const [lenQues, setLenQues] = useState(0) 
    const [question, setQuestion] = useState()
    const [idx, setIdx] = useState(0)
    const [selectedOption, setSelectedOption] = useState(0);
    const [answer, setAnswer] = useState({})
    const [completed, setCompleted] = useState(false)

    // REST calls
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
    useEffect(()=>{
        getQuestionSet()
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

        submitAnswersCall(jsonPayload)
    }

    return(
        <div className="Question Block">
            <h2>{`Welcome ${playerId} in Room: ${roomId}`}</h2>
            {questionSet.Question !== "" && (
                <>
                    <div className="Question">{questionSet.Id}. {questionSet.Question}</div>
                    <ul>
                        <QuestionOptions qs={questionSet} si={selectedOption} ssi={setSelectedOption}/>
                    </ul>
                    {completed ? <button onClick={submit}>Submit</button>  : <button onClick={nxt} id={questionSet.Id}>NextQuestion</button>}
                </>
            )}
        </div>
    )
}

export default Question