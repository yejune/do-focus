import type { Observation } from '../api/client'

interface TimelineProps {
  items: Observation[]
  loading?: boolean
}

export default function Timeline({ items, loading }: TimelineProps) {
  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="animate-pulse flex gap-4">
            <div className="w-2 h-2 bg-gray-300 rounded-full mt-2" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-gray-200 rounded w-1/4" />
              <div className="h-3 bg-gray-200 rounded w-3/4" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (items.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        No observations yet
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {items.map((item) => (
        <div key={item.id} className="flex gap-4">
          <div className="relative">
            <div className={`w-2 h-2 rounded-full mt-2 ${getTypeColor(item.type)}`} />
            <div className="absolute top-4 left-1/2 -translate-x-1/2 w-px h-full bg-gray-200" />
          </div>
          <div className="flex-1 pb-4">
            <div className="flex items-center gap-2 mb-1">
              <span className={`text-xs px-2 py-0.5 rounded ${getTypeBadge(item.type)}`}>
                {item.type}
              </span>
              <span className="text-xs text-gray-500">
                {formatTime(item.created_at)}
              </span>
            </div>
            <p className="text-sm text-gray-700 whitespace-pre-wrap">
              {item.content}
            </p>
          </div>
        </div>
      ))}
    </div>
  )
}

function getTypeColor(type: string): string {
  const colors: Record<string, string> = {
    decision: 'bg-blue-500',
    error: 'bg-red-500',
    success: 'bg-green-500',
    insight: 'bg-purple-500',
    question: 'bg-yellow-500',
  }
  return colors[type] || 'bg-gray-400'
}

function getTypeBadge(type: string): string {
  const badges: Record<string, string> = {
    decision: 'bg-blue-100 text-blue-700',
    error: 'bg-red-100 text-red-700',
    success: 'bg-green-100 text-green-700',
    insight: 'bg-purple-100 text-purple-700',
    question: 'bg-yellow-100 text-yellow-700',
  }
  return badges[type] || 'bg-gray-100 text-gray-700'
}

function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('ko-KR', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
