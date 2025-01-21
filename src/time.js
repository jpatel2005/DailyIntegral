const zeroPadding = (e) => (e < 10 ? '0' + e : e)

const militaryHourToStandard = function (hour) {
  return hour === 0 ? 12 : hour > 12 ? hour - 12 : hour
}

const getTimezone = () =>
  new Date().toLocaleDateString(undefined, { day: '2-digit', timeZoneName: 'short' }).substring(4)

export const convertDateToUTC = function (date) {
  return new Date(
    date.getUTCFullYear(),
    date.getUTCMonth(),
    date.getUTCDate(),
    date.getUTCHours(),
    date.getUTCMinutes(),
    date.getUTCSeconds(),
  )
}

const getUTCMidnight = function () {
  const date = new Date()
  date.setUTCHours(0, 0, 0, 0)
  return date
}

export const getRefreshTime = function () {
  const utcMidnight = getUTCMidnight()
  const localTime = new Date(utcMidnight)
  const hours = localTime.getHours()
  const minutes = localTime.getMinutes()
  const meridian = hours >= 12 ? 'pm' : 'AM'
  return `${militaryHourToStandard(hours)}${minutes === 0 ? '' : `:${zeroPadding(minutes)}`}${meridian} ${getTimezone()}`.trim()
}

export const observeDateRange = function (startDate, endDate) {
  // Get the number of days of the longest streak (use floor just in case)
  const days = Math.floor((endDate.getTime() - startDate.getTime()) / 86400000) + 1
  // Construct longest streak string (acconut for current day problem being solved or not)
  const range = `(${formatDateUS(startDate)} - ${formatDateUS(endDate)})`
  return [days, range]
}

export const formatDateISO = function (date) {
  // Format date in YYYY-MM-DD
  return `${date.getUTCFullYear()}-${zeroPadding(date.getUTCMonth() + 1)}-${zeroPadding(date.getUTCDate())}`
}

export const formatDateUS = function (date) {
  // Format date in MM/YY/DD
  return `${zeroPadding(date.getUTCMonth() + 1)}/${zeroPadding(date.getUTCDate())}/${zeroPadding(date.getUTCFullYear() % 100)}`.trim()
}

const monthNames = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec',
]

export const formatDateString = function (date) {
  const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
  // If date is in the future, then route to home page (prevent trolling)
  const month = date.getUTCMonth()
  const day = date.getUTCDay()
  const monthDate = date.getUTCDate()
  const year = date.getUTCFullYear()
  return `${dayNames[day]}, ${monthNames[month]} ${monthDate < 10 ? '0' + monthDate : monthDate}, ${year}`
}

export const reformatDate = function (dateStr) {
  // Input format: YYYY-MM-DD
  // Output format: <month> <day>, <year>
  if (!/^\d{4}-\d{2}-\d{2}$/.test(dateStr)) {
    throw new Error('Invalid format provided. Expected format is YYYY-MM-DD')
  }
  const [year, month, day] = dateStr.split('-')
  return `${monthNames[+month - 1]} ${+day}, ${year}`
}

export const sleep = (ms) => new Promise((r) => setTimeout(r, ms))
