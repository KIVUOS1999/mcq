import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { connect } from "nats.ws";
import "./Submit.css"

function Submit() {
	const { roomId } = useParams();

	const [parsedJson, setParsedJson] = useState(null);
	const [messages, setMessages] = useState("");
	const [natsSubs, setNatsSubs] = useState();
	const [details, setDetails] = useState("")

	useEffect(() => {
		const connectToNats = async () => {
			const nc = await connect({
				servers: process.env.REACT_APP_NATS_WEBSOCKET,
			});
			const sub = nc.subscribe(roomId);

			setNatsSubs(sub);

			for await (const msg of sub) {
				setMessages(msg);
			}
		};

		connectToNats();
	}, [roomId]);

	// nats player addition and starting game
	useEffect(() => {
		const data = messages._rdata;
		const natsData = new TextDecoder().decode(data);

		// console.log(natsData);

		if (natsData && natsData !== "") {
			try {
				var parsed_json = JSON.parse(natsData)
				setParsedJson(parsed_json);
			} catch {
				// console.log("not a json message", parsed_json)
			}
		}
	}, [messages]);

	// api call to ask server to pusblish result
	useEffect(() => {
		const url = process.env.REACT_APP_BACKEND_BASE + "/get_result/" + roomId;
		fetch(url)
	}, [natsSubs]);

	const getAnswerSet = (val) => {
		return (
			<div className={`detailed_answer_section ${String(val.selected_option) === "0"?"no_answer": (String(val.selected_option) === String(val.correct_option)?"correct_answer":"wrong_answer")}`}>
				<div className="question">{val.question}</div>
				<div className="selected_option">{String(val.selected_option) === "0" ? "Not answered": `Selected Option: ${val.selected_option}`}</div>
				<div className="correct_option">Correct Option: {val.correct_option} - {val.correct_answer}</div>
			</div>
		)
	}

	const getDetails = (player_id)=>{
		console.log(`showing ans for player ${player_id}`)
		console.debug(parsedJson)
		if (player_id == "") return null
		const player_detail = parsedJson.detail.find((item) => item.player_id === player_id);
		console.log("pd", player_detail)
		if (!player_detail) return null
		return Object.entries(player_detail.detailed_answer_set).map(([key, val]) =>(getAnswerSet(val))) 
	}

	const show_details = (player_id) => {
		setDetails(player_id)
	}

	return (
		<div className="scoreCard">
			<div className="scorecard-head">Score</div>
			{parsedJson && (
				<div className="score_card_section">
					<div className="player_scores">
						{parsedJson.score_card.map((item, index) => {
							return (
								<div key={index} className={`${item.correct_answer === -1?"score_element_disable":"score_element"}`} onClick={()=>{show_details(item.player_id)}}>
									{item.player_id} : {item.correct_answer === -1 ? "Yet to answer" : item.correct_answer}
								</div>
							);
						})}
					</div>

					<div className="detail">
						{getDetails(details)}
					</div>
				</div>
			)}
		</div>
	);
}

export default Submit;
