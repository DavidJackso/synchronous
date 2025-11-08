import { Card, Button, Typography } from 'antd';
import { RobotOutlined, ThunderboltOutlined, LockOutlined } from '@ant-design/icons';
import './AIReportTeaser.css';

const { Title, Text, Paragraph } = Typography;

interface AIReportTeaserProps {
  onUpgrade?: () => void;
}

/**
 * AI Report Teaser Component
 * Promotes Pro features with AI-powered analytics
 */
export const AIReportTeaser = ({ onUpgrade }: AIReportTeaserProps) => {
  return (
    <Card className="ai-report-teaser">
      <div className="ai-report-teaser__icon-wrapper">
        <RobotOutlined className="ai-report-teaser__icon" />
        <ThunderboltOutlined className="ai-report-teaser__badge" />
      </div>

      <div className="ai-report-teaser__content">
        <div className="ai-report-teaser__header">
          <LockOutlined className="ai-report-teaser__lock" />
          <Title level={4} className="ai-report-teaser__title">
            AI Анализ Сессии
          </Title>
        </div>

        <Paragraph className="ai-report-teaser__description">
          Получите детальный анализ вашей продуктивности с помощью
          искусственного интеллекта
        </Paragraph>

        <div className="ai-report-teaser__features">
          <div className="ai-report-teaser__feature">
            <span className="ai-report-teaser__feature-icon">✨</span>
            <Text className="ai-report-teaser__feature-text">
              Персональные рекомендации
            </Text>
          </div>
          <div className="ai-report-teaser__feature">
            <span className="ai-report-teaser__feature-icon">📊</span>
            <Text className="ai-report-teaser__feature-text">
              Анализ паттернов продуктивности
            </Text>
          </div>
          <div className="ai-report-teaser__feature">
            <span className="ai-report-teaser__feature-icon">🎯</span>
            <Text className="ai-report-teaser__feature-text">
              Советы по улучшению фокуса
            </Text>
          </div>
        </div>

        <Button
          type="primary"
          size="large"
          block
          onClick={onUpgrade}
          className="ai-report-teaser__upgrade-btn"
        >
          <ThunderboltOutlined />
          Перейти на Pro
        </Button>

        <Text type="secondary" className="ai-report-teaser__hint">
          Попробуйте 7 дней бесплатно
        </Text>
      </div>
    </Card>
  );
};
