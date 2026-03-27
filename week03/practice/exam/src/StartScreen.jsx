// src/StartScreen.jsx
export default function StartScreen({ onStart }) {
  return (
    <div className="center-card">
      <h1 style={{ marginBottom: '10px', color: '#333', fontSize: '28px' }}>
        React 基础知识测试
      </h1>
      <p style={{ color: '#666', marginBottom: '40px', fontSize: '16px' }}>
        本测试旨在检验你对 React 第一阶段基础知识的掌握情况。准备好了吗？
      </p>
      
      {/* 点击按钮触发父组件传来的 onStart 函数，从而改变 App 中的 currentScreen 状态 */}
      <button className="btn btn-primary" style={{ padding: '12px 40px', fontSize: '18px' }} onClick={onStart}>
        开始考试
      </button>
    </div>
  );
}