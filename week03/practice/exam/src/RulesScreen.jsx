// src/RulesScreen.jsx
export default function RulesScreen({ onContinue, onExit }) {
  return (
    <div className="center-card">
      <h2 style={{ marginBottom: '20px', color: '#333' }}>考试须知</h2>
      
      <div style={{ 
        textAlign: 'left', 
        backgroundColor: '#f8f9fa', 
        padding: '25px', 
        borderRadius: '8px', 
        marginBottom: '30px', 
        width: '80%',
        lineHeight: '1.8',
        color: '#444'
      }}>
        <p>1. 考试过程中请勿刷新页面，否则答题进度将丢失。</p>
        <p>2. 每道题仅有一个正确答案。</p>
        <p>3. 答题过程中顶部会有倒计时，超时可能会自动交卷（进阶功能）。</p>
        <p>4. 最终得分将根据答对的题目数量进行计算。</p>
      </div>

      <div style={{ display: 'flex', gap: '20px' }}>
        {/* 调用退出的回调函数，回到 start 页面 */}
        <button className="btn btn-outline" onClick={onExit}>
          退出考试
        </button>
        {/* 调用继续的回调函数，进入 quiz 页面 */}
        <button className="btn btn-primary" onClick={onContinue}>
          我知道了，继续
        </button>
      </div>
    </div>
  );
}