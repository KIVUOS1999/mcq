import { useState } from "react"

const QuestionOptions = (prop)=>{
    const [selectedOption, setSelectedOption] = useState(null)

    const handleCheckBoxChange = (index) => {
        setSelectedOption(index)
        prop.ssi(index)
    }
    
    return(
        Object.values(prop.qs.Options).map((option, index) => (
            <div 
                className={`question_options ${selectedOption === index ? 'question_options_active' : ''}`}
                key={index} id={index+1} onClick={()=>handleCheckBoxChange(index)}>
                {option}
            </div>
        ))
    )
}

export default QuestionOptions