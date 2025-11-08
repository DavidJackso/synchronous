import { Card, Typography, Progress } from 'antd';
import { CheckCircleOutlined, ClockCircleOutlined, FireOutlined } from '@ant-design/icons';
import './SessionStats.css';

const { Title, Text } = Typography;

interface SessionStatsProps {
  tasksCompleted: number;
  tasksTotal: number;
  focusTime: number; // in minutes
  breakTime: number; // in minutes
  cyclesCompleted: number;
}

/**
 * Session Statistics Widget
 * Displays completion stats with visual progress indicators
 */
export const SessionStats = ({
  tasksCompleted,
  tasksTotal,
  focusTime,
  breakTime,
  cyclesCompleted,
}: SessionStatsProps) => {
  const completionRate = tasksTotal > 0 ? Math.round((tasksCompleted / tasksTotal) * 100) : 0;
  const totalTime = focusTime + breakTime;

  return (
    <Card className="session-stats">
      <div className="session-stats__header">
        <Title level={3} className="session-stats__title">
          Статистика сессии
        </Title>
        <div className="session-stats__completion-badge">
          <CheckCircleOutlined className="session-stats__badge-icon" />
          <Text strong className="session-stats__badge-text">
            Завершено
          </Text>
        </div>
      </div>

      {/* Tasks Progress */}
      <div className="session-stats__card session-stats__card--tasks">
        <div className="session-stats__card-header">
          <CheckCircleOutlined className="session-stats__card-icon" />
          <Text className="session-stats__card-label">Задачи</Text>
        </div>
        <div className="session-stats__card-content">
          <div className="session-stats__big-number">
            {tasksCompleted}<span className="session-stats__divider">/</span>{tasksTotal}
          </div>
          <Progress
            percent={completionRate}
            strokeColor={{
              '0%': '#2ecc71',
              '100%': '#27ae60',
            }}
            showInfo={false}
            strokeWidth={8}
            className="session-stats__progress"
          />
          <Text type="secondary" className="session-stats__percent">
            {completionRate}% выполнено
          </Text>
        </div>
      </div>

      {/* Time Stats */}
      <div className="session-stats__grid">
        <div className="session-stats__card session-stats__card--time">
          <div className="session-stats__card-header">
            <ClockCircleOutlined className="session-stats__card-icon" />
            <Text className="session-stats__card-label">Время фокуса</Text>
          </div>
          <div className="session-stats__card-value">
            {focusTime} <span className="session-stats__unit">мин</span>
          </div>
        </div>

        <div className="session-stats__card session-stats__card--cycles">
          <div className="session-stats__card-header">
            <FireOutlined className="session-stats__card-icon" />
            <Text className="session-stats__card-label">Циклы</Text>
          </div>
          <div className="session-stats__card-value">
            {cyclesCompleted}
          </div>
        </div>
      </div>

      {/* Total Time */}
      <div className="session-stats__total">
        <Text type="secondary">Общее время сессии:</Text>
        <Text strong className="session-stats__total-value">
          {totalTime} минут
        </Text>
      </div>
    </Card>
  );
};
