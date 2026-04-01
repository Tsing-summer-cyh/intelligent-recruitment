import React, { useEffect, useState } from 'react';
import { Card, Spin } from 'antd';
import ReactMarkdown from 'react-markdown';
import request from '../../utils/request';

export default function Summary() {
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    request.get('/api/summary').then((res: any) => {
      // 关键：将 markdown 中的静态图片路径替换为完整的接口路径
      const modifiedContent = res.content.replace(/\]\(assets\//g, '](/api/static/assets/');
      setContent(modifiedContent);
    }).finally(() => setLoading(false));
  }, []);

  return (
    <Card title="学习总结" className="min-h-[500px]">
      {loading ? <Spin /> : (
        <div className="prose max-w-none">
          <ReactMarkdown>{content}</ReactMarkdown>
        </div>
      )}
    </Card>
  );
}