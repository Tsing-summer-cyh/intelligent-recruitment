// src/ResultScreen.jsx
export default function ResultScreen({ score, onRestart, onExit }) {
  // 我们有 3 道题，每题算 33.3 分，凑个整
  const finalScore = Math.round((score / 3) * 100);

  return (
    <div className="center-card">
      <h2 style={{ marginBottom: '10px', color: '#333' }}>考试结束！</h2>
      
      <div style={{ margin: '30px 0', textAlign: 'center' }}>
        <p style={{ fontSize: '18px', color: '#666', marginBottom: '10px' }}>你的最终得分是</p>
        <div style={{ fontSize: '72px', fontWeight: 'bold', color: finalScore >= 60 ? '#7ed321' : '#e53e3e' }}>
          {finalScore}
        </div>
        <p style={{ marginTop: '10px', color: '#888' }}>
          共答对 {score} 道题
        </p>
      </div>

      <div style={{ display: 'flex', gap: '20px' }}>
        <button className="btn btn-outline" onClick={onExit}>
          返回首页
        </button>
        <button className="btn btn-primary" onClick={onRestart}>
          重新考试
        </button>
      </div>
    </div>
  );
}