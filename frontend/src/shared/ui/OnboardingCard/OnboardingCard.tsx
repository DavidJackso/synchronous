import { Card, Typography } from 'antd';
import type { OnboardingStep } from '@/shared/types';
import './OnboardingCard.css';

const { Title, Paragraph } = Typography;

interface OnboardingCardProps {
  step: OnboardingStep;
}

export const OnboardingCard = ({ step }: OnboardingCardProps) => {
  return (
    <Card
      className="onboarding-card"
      style={{ background: step.gradient }}
    >
      <div className="onboarding-card__icon">{step.icon}</div>
      <Title level={3} className="onboarding-card__title">
        {step.title}
      </Title>
      <Paragraph className="onboarding-card__description">
        {step.description}
      </Paragraph>
    </Card>
  );
};
