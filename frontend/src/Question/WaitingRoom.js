import { connect } from "nats.ws";
import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import paths from "../constants/constants";
import { FaUserCircle, FaHome, FaCopy } from 'react-icons/fa';
import './WaitingRoom.css'
import {Toaster, toast} from 'sonner'


const WaitingRoom = (props) => {
	const { roomId, playerId, isAdmin, time, questionNumbers } = useParams();

	const [messages, setMessages] = useState([]);
	const [natsSubs, setNatsSubs] = useState();
	const [playersInLobby, setPlayersInLobby] = useState("");

	const isAddPlayerApiCalled = useRef(false);

	const startGame = "startGame:";
	const navigate = useNavigate();

	// settingup nats consumer
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
	}, [props.roomId]);

	// nats player addition and starting game
	useEffect(() => {
		const data = messages._rdata;
		const natsData = new TextDecoder().decode(data);

		if (natsData.startsWith(startGame)) {
			// console.log("game has started");
			const timeData = natsData.split(":")[1];

			const url = `${process.env.REACT_APP_BACKEND_BASE}${paths.g_mcq}/room/${roomId}/count/10/admin/false`;
			fetch(url)
				.then((res) => {
					if (res.ok) {
						navigate(`/question/${roomId}/${playerId}/${timeData}/${isAdmin}`, {replace: true});
					}
					if (!res.ok) {
						throw new Error(`HTTP fetch error Status: ${res.status}`);
					}
				}).catch((error) => {
					toast.error("error creating question", error)
				})

			return
		}

		setPlayersInLobby(natsData);
	}, [messages]);

	useEffect(() => {
		if (natsSubs !== undefined && !isAddPlayerApiCalled.current) {
			addPlayer();
			isAddPlayerApiCalled.current = true;
		}
	}, [natsSubs]);

	const copyBtnOnClick = () => {
		var room_code = document.getElementById("room_id").innerHTML
		navigator.clipboard.writeText(room_code)
		toast.success("Room id copied")
	}

	const addPlayer = () => {
		const url =
			process.env.REACT_APP_BACKEND_BASE +
			paths.a_player +
			roomId +
			"/player/" +
			playerId +
			"/admin/true";
		fetch(url)
			.then((res) => {
				if (!res.ok) {
					if (res.status === 400) {
						res.json().then((errData) => {
							alert(
								`Error in AddPlayer: ${errData.response_code}|${errData.reason}`,
							);
						});
						return;
					}
					throw new Error(`HTTP: call error: ${res.status}`);
				}
			})
			.catch((error) => {
				console.error(error);
			});
	};

	// admin start button action
	const sendStartMessage = () => {
		const url =
			process.env.REACT_APP_BACKEND_BASE +
			paths.s_game +
			roomId +
			"/endtime/" +
			time * 60 * 1000;
		fetch(url)
			.then((res) => {
				if (!res.ok) {
					if (res.status === 400) {
						console.error("Error in nats start game");
						return;
					}
				}
			})
			.catch((error) => {
				console.error(error);
			});
	};
	
	let playersArray = playersInLobby.split(';').filter(player => player.trim() !== '');

	return (
		<div className="Waiting_room">
			<div className="content">
				<div className="Lobby">
					<div className="Room_icon"><FaHome /></div>
					<div className="Room_id" id="room_id">{roomId}</div>
					<div className="Room_copy" id="copy_btn" onClick={copyBtnOnClick}><FaCopy /></div>
				</div>
				<div className="Total_count"><b>Joined</b> {playersArray.length}</div>
				<div className="Player_id">
					{playersArray.map((player, index) => (
						<div className={`player_element ${player===playerId?"CurrentPlayer":""}`}>
							<div className="player_icon">
								<FaUserCircle />
							</div>
							<div className="player_name" key={index}>
								{player}
							</div>
						</div>
					))}
				</div>
			</div>
			<div className="Start-Button">
				{isAdmin == 1 ? <button className="Start_game" onClick={sendStartMessage}>Start</button> : ""}
			</div>
			<Toaster />
		</div>
	);
};

export default WaitingRoom;
