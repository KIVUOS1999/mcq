import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { connect } from "nats.ws";

function Submit() {
	const { roomId } = useParams();

	const [parsedJson, setParsedJson] = useState(null);
	const [messages, setMessages] = useState("");
	const [natsSubs, setNatsSubs] = useState();

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

	return (
		<div className="ScoreCard">
			{parsedJson && (
				<div className="score_card_section">
					<h1>ScoreCard : {roomId}</h1>
					<div className="Players">
						<ul>
							{parsedJson.score_card.map((item, index) => {
								return (
									<li key={index}>
										{item.player_id} :{" "}
										{item.correct_answer === -1
											? "Yet to answer"
											: item.correct_answer}
									</li>
								);
							})}
						</ul>
						<h2>Details</h2>

						{parsedJson.detail.map((item, index) => {
							return (
								<div className="details">
									<h3 className="player_name">player name: {item.player_id}</h3>
									<div className="detailed_answer">
										{Object.entries(item.detailed_answer_set).map(([key, val], idx) => {
											return (
												<div className="detailed_answer_section">
													<div className="question">Q. {val.question}</div>
													<div className="selected_option">Selected Option: {val.selected_option}</div>
													<div className="correct_option">Correct Option: {val.correct_option} - {val.correct_answer}</div>
													<div className="is_correct">{String(val.selected_option) === String(val.correct_option) ? "Correct": "Incorrect"}</div>
													<br/>
												</div>
											)
										})}
									</div>
								</div>
							)
						})}

					</div>
				</div>
			)}
		</div>
	);
}

export default Submit;
