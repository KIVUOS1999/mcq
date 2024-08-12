import { useState } from "react"
import CreateRoom from "./CreateRoom"
import JoinRoom from "./JoinRoom"
import './Room.css'

function Room() {
    const [createRoom, setCreateRoom] = useState(false)

    return(
        <div className="RoomBody">
            <div className="RoomContainer">
                <div className="ActionButtons">
                    <button className={`Btn CreateRoomButton ${createRoom? 'ActiveBtn':''}`} onClick={()=>{
                        setCreateRoom(true)
                    }} disabled={createRoom}>Create</button>
                    
                    <button className={`Btn JoinRoomButton ${!createRoom? 'ActiveBtn':''}`} onClick={()=>{
                        setCreateRoom(false)
                    }} disabled={!createRoom}>Join</button>
                </div>
                
                <div className="ActionPart">
                    {createRoom ? <CreateRoom/> : <JoinRoom/>}
                </div>
            </div>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320"><path fill="#990011" fill-opacity="1" d="M0,256L15,245.3C30,235,60,213,90,213.3C120,213,150,235,180,256C210,277,240,299,270,288C300,277,330,235,360,229.3C390,224,420,256,450,229.3C480,203,510,117,540,74.7C570,32,600,32,630,64C660,96,690,160,720,176C750,192,780,160,810,144C840,128,870,128,900,117.3C930,107,960,85,990,117.3C1020,149,1050,235,1080,261.3C1110,288,1140,256,1170,256C1200,256,1230,288,1260,261.3C1290,235,1320,149,1350,144C1380,139,1410,213,1425,250.7L1440,288L1440,320L1425,320C1410,320,1380,320,1350,320C1320,320,1290,320,1260,320C1230,320,1200,320,1170,320C1140,320,1110,320,1080,320C1050,320,1020,320,990,320C960,320,930,320,900,320C870,320,840,320,810,320C780,320,750,320,720,320C690,320,660,320,630,320C600,320,570,320,540,320C510,320,480,320,450,320C420,320,390,320,360,320C330,320,300,320,270,320C240,320,210,320,180,320C150,320,120,320,90,320C60,320,30,320,15,320L0,320Z"></path></svg>
        </div>
    )
}

export default Room