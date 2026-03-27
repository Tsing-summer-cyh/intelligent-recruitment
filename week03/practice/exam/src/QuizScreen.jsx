// src/QuizScreen.jsx
import { useState, useEffect } from 'react';

const questions = [
  {
    id: 1,
    question: '在 React 中，用于声明组件状态的 Hook 是什么？',
    options: ['useEffect', 'useState', 'useRef', 'useContext'],
    answer: 'useState'
  },
  {
    id: 2,
    question: 'React 的函数组件名称必须以什么开头？',
    options: ['小写字母', '大写字母', '数字', '特殊符号'],
    answer: '大写字母'
  },
  {
    id: 3,
    question: '以下哪个 Hook 通常用于处理发送网络请求等“副作用”？',
    options: ['useState', 'useEffect', 'useMemo', 'useCallback'],
    answer: 'useEffect'
  }
];

export default function QuizScreen({ onFinish }) {
  const [currentIndex, setCurrentIndex] = useState(0); 
  const [score, setScore] = useState(0); 
  const [timeLeft, setTimeLeft] = useState(60); 
  const [selectedOption, setSelectedOption] = useState(null); 
  const [isAnswered, setIsAnswered] = useState(false); 
  const [transitionTimeoutId, setTransitionTimeoutId] = useState(null);

  useEffect(() => {
    if (timeLeft === 0) {
      onFinish(score);
      return;
    }
    if (isAnswered) return;
    
    const timer = setInterval(() => {
      setTimeLeft(prev => prev - 1);
    }, 1000);

    return () => clearInterval(timer);
  }, [timeLeft, score, onFinish, isAnswered]);

  // ============ 🌟 修复关键点 1：接收最新分数参数 ============
  const handleNextQuestion = (latestScore) => {
    if (transitionTimeoutId) {
      clearTimeout(transitionTimeoutId);
      setTransitionTimeoutId(null);
    }

    // 判断：如果传了最新分数（通过延时器传来的），就用最新的；否则用 state 里的 score（比如用户手动点下一题）
    const finalScore = typeof latestScore === 'number' ? latestScore : score;

    if (currentIndex + 1 < questions.length) {
      setCurrentIndex(prevIndex => prevIndex + 1); 
      setSelectedOption(null); 
      setIsAnswered(false); 
    } else {
      onFinish(finalScore); // 🌟 使用准确的分数交卷
    }
  };

  const handleAnswerClick = (option) => {
    if (isAnswered) return; 

    setIsAnswered(true); 
    setSelectedOption(option); 

    const isCorrect = option === questions[currentIndex].answer;
    const newScore = isCorrect ? score + 1 : score;
    setScore(newScore);

    const timeoutId = setTimeout(() => {
      // ============ 🌟 修复关键点 2：将算好的 newScore 直接传进去 ============
      handleNextQuestion(newScore);
    }, 1500); 

    setTransitionTimeoutId(timeoutId);
  };

  const currentQ = questions[currentIndex];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div className="content-area" style={{ flex: '1 0 auto' }}>
        <div className="header" style={{ padding: '0 0 15px 0', borderBottom: 'none' }}>
          <span style={{ fontSize: '18px', fontWeight: 'bold', color: '#333' }}>
            题目 {currentIndex + 1} / {questions.length}
          </span>
          <span style={{ color: '#e53e3e', fontWeight: 'bold' }}>
            倒计时: {timeLeft}s
          </span>
        </div>
        
        <div style={{ marginTop: '20px' }}>
          <h3 style={{ marginBottom: '30px', color: '#2d3748', lineHeight: '1.5' }}>
            {currentQ.question}
          </h3>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {currentQ.options.map((option, index) => {
              let optionClassName = "option-btn";
              if (isAnswered) {
                const isCorrectAnswer = option === currentQ.answer;
                const isUserSelection = option === selectedOption;
                if (isCorrectAnswer) { optionClassName += " correct"; } 
                else if (isUserSelection && !isCorrectAnswer) { optionClassName += " wrong"; }
              }

              return (
                <button 
                  key={index}
                  className={optionClassName}
                  onClick={() => handleAnswerClick(option)}
                  disabled={isAnswered} 
                >
                  {option}
                  {isAnswered && (option === currentQ.answer ? '✔️' : (option === selectedOption ? '✖️' : ''))}
                </button>
              );
            })}
          </div>
        </div>
      </div>

      <div className="quiz-footer">
        <div className="progress-dots">
          {questions.map((_, index) => (
            <div 
              key={index}
              className={`dot ${index === currentIndex ? 'current' : ''}`} 
            />
          ))}
        </div>

        <button 
          className="btn btn-primary btn-next"
          // ============ 🌟 修复关键点 3：改写为箭头函数，防止事件对象 e 被错误传入 ============
          onClick={() => handleNextQuestion()} 
        >
          {currentIndex + 1 < questions.length ? '下一题 ->' : '结束考试 ->'}
        </button>
      </div>
    </div>
  );
}