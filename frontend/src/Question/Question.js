import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { connect } from "nats.ws";

import paths from "../constants/constants";
import QuestionOptions from "./QuestionOptions";
import "./Question.css"

import { GrCaretNext } from "react-icons/gr";
import { FaUserAstronaut } from "react-icons/fa";
import { toast, Toaster } from "sonner";

function Question() {
	const { roomId, playerId, time, admin } = useParams();
	const navigate = useNavigate();
	const [questionSet, setQuestionSet] = useState({
		Question: "",
		Options: [],
		Id: "",
	});

	const [lenQues, setLenQues] = useState(0);
	const [question, setQuestion] = useState();
	const [idx, setIdx] = useState(0);
	const [selectedOption, setSelectedOption] = useState(-1);
	const [questionID, setQuestionID] = useState("")
	const [answer, setAnswer] = useState({});
	const [completed, setCompleted] = useState(false);
	const [countdownTimer, setContDownTimer] = useState(time);

	const ns = useRef(false);

	const [messages, setMessages] = useState();

	// nats server
	const connectToNats = async () => {
		// console.log("connecting to nats server from question", roomId);

		const nc = await connect({ servers: process.env.REACT_APP_NATS_WEBSOCKET });
		const sub = nc.subscribe(roomId);

		for await (const msg of sub) {
			setMessages(msg);
		}

		ns.current = true;
	};

	// getting questions set from server
	const getQuestionSet = (route) => {
		const url = `${process.env.REACT_APP_BACKEND_BASE}${route}`;

		fetch(url)
			.then((res) => {
				if (!res.ok) {
					throw new Error(`HTTP fetch error Status: ${res.status}`);
				}

				return res.json();
			})
			.then((data) => {
				// console.log("fetched question set:", data)
				setQuestion(data);
				setLenQues(data.data.length);
			})
			.catch((error) => {
				console.error(`Error fetching data, ${error}`);
			});
	};

	// submitting answers for the player
	const submitAnswersCall = (jsonBody) => {
		console.debug("json", jsonBody)
		const url = `${process.env.REACT_APP_BACKEND_BASE}${paths.a_answer}`;
		fetch(url, {
			method: "POST",
			body: jsonBody,
		})
			.then((res) => {
				if (res.status !== 201) {
					if (res.status === 400) {
						res.json().then((errData) => {
							alert(
								`Error in Submitting answers: ${errData.response_code}|${errData.reason}`,
							);
						});
						return;
					}
					throw new Error(`HTTP fetch error Status: ${res.status}`);
				} else {
					navigate(`/submit/${roomId}`, {replace: true});
				}
			})
			.catch((error) => {
				console.error(`Error fetching data, ${error}`);
			});
	};

	// Use Effects
	useEffect(() => {
		if (
			messages &&
			new TextDecoder().decode(messages._rdata) === "submit_game"
		) {
			submit();
		}
	}, [messages]);

	useEffect(() => {
		getQuestionSet(`${paths.g_mcq}/room/${roomId}/count/10/admin/false`)

		if (ns.current == false) {
			connectToNats();
		}

		var x = setInterval(() => {
			setContDownTimer((prevTime) => prevTime - 1);
		}, 1000);

		return () => clearInterval(x);
		// eslint-disable-next-line
	}, []);

	useEffect(() => {
		if (question && question.data[idx]) {
			setQuestionID(question.data[idx].id)
			setQuestionSet({
				Question: question.data[idx].question,
				Options: question.data[idx].options,
				Id: question.data[idx].id,
			});
		}
	}, [idx, question]);

	useEffect(() => {
		if (Object.keys(answer).length === 0) {
			return;
		}
		submit();
	}, [completed]);

	// Button Functions
	let nxt = (e) => {
		if(selectedOption == -1) {
			toast.error("Select an option")

			return
		}

		const questionIndex = questionID;
		const so = selectedOption + 1;

		setAnswer({ ...answer, [questionIndex]: so });

		setIdx(idx + 1);

		if (idx >= lenQues) {
			setCompleted(true)
		}
	};

	let submit = () => {
		let preparedJson = {
			room_id: roomId,
			player_id: playerId,
			answer_set: answer,
		};

		const jsonPayload = JSON.stringify(preparedJson);

		// console.log(jsonPayload);
		submitAnswersCall(jsonPayload);
	};

	return (
		<div className="question_block">
			<div className="question_header">
				<div className="userID"><FaUserAstronaut/> {playerId}</div>
				<div className="timer">{countdownTimer} sec</div>
			</div>
			{questionSet.Question !== "" && (
				<div className="question_container">
					<div className="question_section">
						<div className="question_question">
							Q: {questionSet.Question}
						</div>
						
						<div className="question_options_conatiner">
							<QuestionOptions
								qs={questionSet}
								si={selectedOption}
								ssi={setSelectedOption}
							/>
						</div>
					</div>
					
					<button className="question_submit" onClick={nxt} id={questionSet.Id}>
						<GrCaretNext />
					</button>
					
					
				</div>
			)}
			<Toaster />
		</div>
	);
}

export default Question;
