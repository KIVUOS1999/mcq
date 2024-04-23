const QuestionOptions = (prop)=>{
    const handleCheckBoxChange = (index) => {
        prop.ssi(index)
    }
    
    return(
        Object.values(prop.qs.Options).map((option, index) => (
            <li key={index}>
                <input 
                type="radio" 
                id={index} 
                name={option} 
                value={index+1} 
                checked={index === prop.si}
                onChange={()=>handleCheckBoxChange(index)}/>
                <label htmlFor="{index}">{option}</label>
            </li>
        ))
    )
}

export default QuestionOptions