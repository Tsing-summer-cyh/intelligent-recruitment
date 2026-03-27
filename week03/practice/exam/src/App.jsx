// src/App.jsx
import { useState } from 'react';
import StartScreen from './StartScreen';
import RulesScreen from './RulesScreen';
import QuizScreen from './QuizScreen';
import ResultScreen from './ResultScreen';

export default function App() {
  // 核心状态：当前处于哪个页面 ('start' | 'rules' | 'quiz' | 'result')
  const [currentScreen, setCurrentScreen] = useState('start');
  // 核心状态：总得分
  const [score, setScore] = useState(0);

  // 页面切换的渲染逻辑
  const renderScreen = () => {
    switch (currentScreen) {
      case 'start':
        // 传入切换到规则页的方法
        return <StartScreen onStart={() => setCurrentScreen('rules')} />;
      case 'rules':
        // 传入继续考试和退出考试的方法
        return (
          <RulesScreen 
            onContinue={() => setCurrentScreen('quiz')} 
            onExit={() => setCurrentScreen('start')} 
          />
        );
      case 'quiz':
        // 传入考试结束后的回调，接收分数并跳转到结果页
        return (
          <QuizScreen 
            onFinish={(finalScore) => {
              setScore(finalScore);
              setCurrentScreen('result');
            }} 
          />
        );
      case 'result':
        // 传入得分，以及重新开始/退出的方法
        return (
          <ResultScreen 
            score={score} 
            onRestart={() => {
              setScore(0);
              setCurrentScreen('quiz');
            }}
            onExit={() => setCurrentScreen('start')}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div className="app-container">
      {renderScreen()}
    </div>
  );
}