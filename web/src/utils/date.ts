/**
 * 日期格式化工具函数
 */

/**
 * 格式化日期时间
 * @param date 日期字符串或Date对象
 * @param format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的日期字符串
 */
export function formatDateTime(date: string | Date | null | undefined, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!date) return '-'
  
  const d = typeof date === 'string' ? new Date(date) : date
  
  if (!(d instanceof Date) || isNaN(d.getTime())) {
    return '-'
  }
  
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')
  
  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化日期
 * @param date 日期字符串或Date对象
 * @returns 格式化后的日期字符串 (YYYY-MM-DD)
 */
export function formatDate(date: string | Date | null | undefined): string {
  return formatDateTime(date, 'YYYY-MM-DD')
}

/**
 * 格式化时间
 * @param date 日期字符串或Date对象
 * @returns 格式化后的时间字符串 (HH:mm:ss)
 */
export function formatTime(date: string | Date | null | undefined): string {
  return formatDateTime(date, 'HH:mm:ss')
}

/**
 * 获取相对时间
 * @param date 日期字符串或Date对象
 * @returns 相对时间字符串
 */
export function getRelativeTime(date: string | Date): string {
  if (!date) return '-'
  
  const d = typeof date === 'string' ? new Date(date) : date
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(months / 12)
  
  if (years > 0) return `${years}年前`
  if (months > 0) return `${months}个月前`
  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  if (seconds > 0) return `${seconds}秒前`
  
  return '刚刚'
}

/**
 * 判断是否为今天
 * @param date 日期字符串或Date对象
 * @returns 是否为今天
 */
export function isToday(date: string | Date): boolean {
  if (!date) return false
  
  const d = typeof date === 'string' ? new Date(date) : date
  const today = new Date()
  
  return d.toDateString() === today.toDateString()
}

/**
 * 判断是否为本周
 * @param date 日期字符串或Date对象
 * @returns 是否为本周
 */
export function isThisWeek(date: string | Date): boolean {
  if (!date) return false
  
  const d = typeof date === 'string' ? new Date(date) : date
  const today = new Date()
  const startOfWeek = new Date(today.setDate(today.getDate() - today.getDay()))
  const endOfWeek = new Date(today.setDate(today.getDate() - today.getDay() + 6))
  
  return d >= startOfWeek && d <= endOfWeek
}

/**
 * 获取日期范围文本
 * @param startDate 开始日期
 * @param endDate 结束日期
 * @returns 日期范围文本
 */
export function getDateRangeText(startDate: string | Date, endDate: string | Date): string {
  const start = formatDate(startDate)
  const end = formatDate(endDate)
  
  if (start === '-' && end === '-') return '-'
  if (start === '-') return `至 ${end}`
  if (end === '-') return `${start} 至今`
  if (start === end) return start
  
  return `${start} 至 ${end}`
}