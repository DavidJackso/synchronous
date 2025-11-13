import { Avatar, Progress, Spin } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import type { ParticipantProgress } from '@/shared/api';
import './ParticipantsProgress.css';

interface ParticipantsProgressProps {
  participants: ParticipantProgress[];
  isLoading?: boolean;
}

/**
 * ParticipantsProgress component
 * Displays progress of all participants in a group session
 */
export const ParticipantsProgress: React.FC<ParticipantsProgressProps> = ({
  participants,
  isLoading = false,
}) => {
  if (isLoading) {
    return (
      <div className="participants-progress participants-progress--loading">
        <Spin size="small" />
      </div>
    );
  }

  if (participants.length === 0) {
    return null;
  }

  return (
    <div className="participants-progress">
      <div className="participants-progress__title">
        Прогресс участников
      </div>
      
      <div className="participants-progress__list">
        {participants.map((participant) => (
          <div key={participant.userId} className="participants-progress__item">
            <div className="participants-progress__user">
              <Avatar
                src={participant.avatarUrl}
                icon={<UserOutlined />}
                size={32}
                className="participants-progress__avatar"
              />
              <div className="participants-progress__info">
                <span className="participants-progress__name">
                  {participant.userName}
                </span>
                <span className="participants-progress__stats">
                  {participant.tasksCompleted} / {participant.tasksTotal} задач
                </span>
              </div>
            </div>
            
            <div className="participants-progress__bar">
              <Progress
                percent={Math.round(participant.progressPercent)}
                size="small"
                strokeColor={{
                  '0%': 'var(--color-primary)',
                  '100%': 'var(--color-success)',
                }}
                trailColor="rgba(255, 255, 255, 0.08)"
                showInfo={true}
                format={(percent) => `${percent}%`}
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
