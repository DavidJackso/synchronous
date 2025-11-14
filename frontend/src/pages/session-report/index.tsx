import { useState, useEffect } from 'react';
import { Button, Spin, message } from 'antd';
import { useNavigate, useParams } from 'react-router';
import { HomeOutlined } from '@ant-design/icons';
import { useAppSelector, useAppDispatch } from '@/shared/hooks/redux';
import {
  selectIsGroupMode,
  selectSessionTasks,
  selectCurrentCycle,
} from '@/entities/session/model/activeSessionSelectors';
import { resetSessionSetup } from '@/entities/session/model/sessionSetupSlice';
import { SessionStats } from '@/widgets/session-stats/ui';
import { Leaderboard } from '@/widgets/leaderboard/ui';
import { AIReportTeaser } from '@/features/ai-assistant/ui';
import { sessionsApi, leaderboardApi } from '@/shared/api';
import { useMaxWebApp } from '@/shared/hooks/useMaxWebApp';
import type { Task } from '@/shared/types';
import type { SessionReport, LeaderboardEntry as ApiLeaderboardEntry } from '@/shared/api';
import './styles.css';

// Leaderboard component expects different type
interface ComponentLeaderboardEntry {
  user: {
    id: string;
    name: string;
    avatar: string;
  };
  tasksCompleted: number;
  focusTime: number;
  score: number;
}

/**
 * Session Report Page
 * Shows session results with stats, leaderboard, и ключевые выводы
 */
export function SessionReportPage() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { sessionId } = useParams<{ sessionId: string }>();
  const { isMaxEnvironment, isReady } = useMaxWebApp();
  const isGroupMode = useAppSelector(selectIsGroupMode);
  const localTasks = useAppSelector(selectSessionTasks) as Task[];
  const currentCycle = useAppSelector(selectCurrentCycle);

  const [report, setReport] = useState<SessionReport | null>(null);
  const [leaderboard, setLeaderboard] = useState<ApiLeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Load session report data
  useEffect(() => {
    if (!sessionId) {
      navigate('/');
      return;
    }

    // Wait for MAX WebApp to initialize
    if (!isReady) {
      return;
    }

    const loadReportData = async () => {
      // Dev mode: skip API calls, use local Redux state
      if (!isMaxEnvironment) {
        setIsLoading(false);
        return;
      }

      // Production: load report data from API
      try {
        // Сначала получаем информацию о сессии, чтобы проверить статус
        const sessionResponse = await sessionsApi.getSessionById(sessionId);
        const session = sessionResponse.session;
        
        // Если сессия уже завершена, получаем отчет из данных сессии
        // Если нет - завершаем сессию
        let reportData: SessionReport;
        if (session.status === 'completed') {
          // Сессия уже завершена, формируем отчет из данных сессии
          const sessionTasks = session.tasks ?? [];
          const completedTasks = sessionTasks.filter((t) => t.completed).length;
          const cycles = currentCycle || 1; // Используем currentCycle из Redux
          reportData = {
            sessionId: session.id,
            tasksCompleted: completedTasks,
            tasksTotal: sessionTasks.length,
            focusTime: session.focusDuration * cycles,
            breakTime: session.breakDuration * cycles,
            cyclesCompleted: cycles,
            completedAt: session.completedAt || new Date().toISOString(),
            participants: session.participants.map(p => ({
              userId: p.userId,
              userName: p.userName,
              avatarUrl: p.avatarUrl,
              tasksCompleted: 0, // TODO: получить из статистики
              focusTime: 0, // TODO: получить из статистики
            })),
          };
        } else {
          // Сессия еще не завершена, завершаем её
          const reportResponse = await sessionsApi.completeSession(sessionId);
          reportData = reportResponse.report;
        }
        
        setReport(reportData);

        // Load session leaderboard (for group sessions)
        if (isGroupMode) {
          const leaderboardResponse = await leaderboardApi.getSessionLeaderboard(sessionId);
          setLeaderboard(leaderboardResponse.leaderboard);
        }
      } catch (error) {
        console.error('[SessionReport] Failed to load report:', error);
        message.error('Не удалось загрузить отчёт сессии');
      } finally {
        setIsLoading(false);
      }
    };

    loadReportData();
  }, [sessionId, isGroupMode, isMaxEnvironment, isReady, navigate]);

  // Calculate stats from report or fallback to local data
  const fallbackCompletedTasks = localTasks.filter((t) => t.completed).length;
  const tasksCompleted = report?.tasksCompleted ?? fallbackCompletedTasks;
  const tasksTotal = report?.tasksTotal ?? localTasks.length;
  const cyclesCompleted = report?.cyclesCompleted ?? currentCycle ?? 1;
  const focusTime = report?.focusTime ?? cyclesCompleted * 25;
  const breakTime = report?.breakTime ?? cyclesCompleted * 5;

  // Transform API leaderboard entries to component format
  const leaderboardEntries: ComponentLeaderboardEntry[] = (
    report?.participants || leaderboard
  )
    .map((p) => ({
      user: {
        id: p.userId,
        name: p.userName,
        avatar: p.avatarUrl,
      },
      tasksCompleted: p.tasksCompleted ?? 0,
      focusTime: p.focusTime ?? 0,
      score:
        'score' in p
          ? (p as ApiLeaderboardEntry).score
          : (p.tasksCompleted ?? 0) * 100 + (p.focusTime ?? 0),
    }))
    .sort((a, b) => b.score - a.score || b.tasksCompleted - a.tasksCompleted);

  const handleGoHome = () => {
    dispatch(resetSessionSetup());
    navigate('/');
  };

  const handleUpgrade = () => {
    console.log('[SessionReport] Upgrade to Pro clicked');
    // TODO: Implement upgrade flow
  };

  if (isLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
        <Spin size="large" />
      </div>
    );
  }

  const completionRate = tasksTotal > 0 ? Math.round((tasksCompleted / tasksTotal) * 100) : 0;
  const formattedCompletionDate = report?.completedAt
    ? new Date(report.completedAt).toLocaleString('ru-RU', {
        day: 'numeric',
        month: 'long',
        hour: '2-digit',
        minute: '2-digit',
      })
    : null;

  const heroSummary = isGroupMode
    ? `Команда закрыла ${completionRate}% плановых задач и сфокусирована ${focusTime} мин`
    : `Вы закрыли ${completionRate}% задач и сфокусированы ${focusTime} мин`;

  const insights: string[] = [
    tasksTotal > 0 ? `Закрыто ${tasksCompleted} из ${tasksTotal} задач` : 'Задачи для сессии не заданы',
    `Фокус ${focusTime} мин • Перерывы ${breakTime} мин`,
    `Циклов завершено: ${cyclesCompleted}`,
  ];

  if (isGroupMode && leaderboardEntries.length > 0) {
    const leader = leaderboardEntries[0];
    insights.push(`${leader.user.name} лидирует с результатом ${leader.score} баллов`);
  }

  return (
    <div className="session-report-page">
      <div className="session-report-page__container">
        <section className="session-report-page__hero">
          <div className="session-report-page__hero-info">
            <span className="session-report-page__status">Сессия завершена</span>
            <h1>Отчёт по сессии</h1>
            <p className="session-report-page__summary">{heroSummary}</p>
            {formattedCompletionDate && (
              <span className="session-report-page__meta">Закончено {formattedCompletionDate}</span>
            )}
          </div>
          <div className="session-report-page__hero-stats">
            <div className="session-report-page__hero-stat">
              <span className="session-report-page__hero-label">Выполнение</span>
              <span className="session-report-page__hero-value">{completionRate}%</span>
            </div>
            <div className="session-report-page__hero-stat">
              <span className="session-report-page__hero-label">Задачи</span>
              <span className="session-report-page__hero-value">
                {tasksCompleted}/{tasksTotal || '—'}
              </span>
            </div>
            <div className="session-report-page__hero-stat">
              <span className="session-report-page__hero-label">Циклы</span>
              <span className="session-report-page__hero-value">{cyclesCompleted}</span>
            </div>
          </div>
        </section>

        <div className="session-report-page__grid">
          <SessionStats
            tasksCompleted={tasksCompleted}
            tasksTotal={tasksTotal}
            focusTime={focusTime}
            breakTime={breakTime}
            cyclesCompleted={cyclesCompleted}
          />

          {isGroupMode && leaderboardEntries.length > 0 && (
            <Leaderboard entries={leaderboardEntries} />
          )}
        </div>

        <section className="session-report-page__insights">
          <h3>Ключевые выводы</h3>
          <ul className="session-report-page__insights-list">
            {insights.map((insight) => (
              <li key={insight}>{insight}</li>
            ))}
          </ul>
        </section>

        <AIReportTeaser onUpgrade={handleUpgrade} />

        {/* Home Button */}
        <Button
          type="default"
          size="large"
          block
          icon={<HomeOutlined />}
          onClick={handleGoHome}
          className="session-report-page__home-btn"
        >
          Вернуться на главную
        </Button>
      </div>
    </div>
  );
}
