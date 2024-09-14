import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import paths from "../constants/constants";
import { toast, Toaster } from "sonner";

function isValidInput(value) {
    const isOK = /^[a-zA-Z0-9]+$/.test(value); // Regex for letters and numbers only
    console.log(isOK)

    return isOK
}

function CreateRoom() {
	const [playerID, setPlayerID] = useState("");
	const [roomID, setRoomID] = useState("");
	const [time, setTime] = useState("");
	const [questionNumbers, setQuestionNumbers] = useState(0)

	const navigate = useNavigate();

	const navigateToWaitingRoom = () => {
		navigate(`/lobby/${roomID}/${playerID}/1/${time}`, {replace: true});
	};

	const createRoomAPI = () => {
		const url = process.env.REACT_APP_BACKEND_BASE + paths.c_room;
		fetch(url)
			.then((res) => {
				if (!res.ok) {
					if (res.status === 400) {
						res.json().then((errData) => {
							alert(
								`Error in CreateRoom: ${errData.response_code}|${errData.reason}`,
							);
						});
						return;
					}
					throw new Error(`HTTP: call error: ${res.status}`);
				}

				return res.json();
			})
			.then((data) => {
				const roomNo = data.room_id;
				setRoomID(roomNo);

				const player = document.getElementById("player_name");
				const player_time = document.getElementById("time").innerHTML;

				setPlayerID(player.value);
				setTime(player_time);
			})
			.catch((error) => {
				toast.error(error)
			});
	};

	const createRoomAddPlayer = () => {
		const player = document.getElementById("player_name")
		const player_time = document.getElementById("time").innerHTML
		console.log("pt",player_time)
		const question_nos = document.getElementById("questions").value

		if (!isValidInput(player.value)){
            toast.error("Make your name cool but use numbers and alphabets")
            return
        }

		if (question_nos <= 0 ) {
			toast.error("you don't require questions?")
			return
		}

		createRoomAPI(player.value, player_time, question_nos);
	};

	useEffect(() => {
		if (roomID !== "" && playerID !== "") {
			navigateToWaitingRoom();
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [roomID, playerID]);

	// Event listener for time increase
	let intervalId = null
	const intervalTime = 250

	useEffect(() => {
		const increaseButton = document.getElementById("TimeInc");
		increaseButton.addEventListener("mousedown", startIncrementing);
		increaseButton.addEventListener("mouseup", stopIncrementing);
		increaseButton.addEventListener("mouseleave", stopIncrementing);
		increaseButton.addEventListener("click", incTime)

		const decreaseButton = document.getElementById("TimeDec")
		decreaseButton.addEventListener("mousedown", startDecrementing);
		decreaseButton.addEventListener("mouseup", stopIncrementing);
		decreaseButton.addEventListener("mouseleave", stopIncrementing);
		decreaseButton.addEventListener("click", decTime)
	}, [])

	const incTime = () => {
		const timeElement = document.getElementById("time");
		let currentTime = parseInt(timeElement.innerHTML, 10);
		if (currentTime >= 60){
			return
		}
		currentTime++;
		timeElement.innerHTML = currentTime;
	};

	const decTime = () => {
		const timeElement = document.getElementById("time");
		let currentTime = parseInt(timeElement.innerHTML, 10);
		if (currentTime <= 1) {
			return
		}
		currentTime--;
		timeElement.innerHTML = currentTime;
	};

	const startIncrementing = ()=> {
		if (intervalId !== null) {
			clearInterval(intervalId)
		}

		intervalId = setInterval(incTime, intervalTime)
	}

	const startDecrementing = ()=> {
		if (intervalId !== null) {
			clearInterval(intervalId)
		}

		intervalId = setInterval(decTime, intervalTime)
	}

	const stopIncrementing = ()=>{
		if (intervalId !== null) {
			clearInterval(intervalId);
			intervalId = null;
		}
	}

	return (
		<div className="CreateRoom">
			<input className="Text Name" id="player_name" type="text" placeholder="Name" required></input>
			<input className="Text Question" id="questions" type="number" placeholder="Questions" required></input>
			<div className="Time">
				<div className="TimeActions">
					<button className="Btn TimeBtn" id="TimeDec">-</button>
					<div id="time">1</div><span>minute</span>
					<button className="Btn TimeBtn" id="TimeInc">+</button>
				</div>
			</div>
			<button className="Btn ActionBtn" onClick={createRoomAddPlayer}>Create Room</button>
			<Toaster/>
		</div>
	);
}

export default CreateRoom;
