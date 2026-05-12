// Shared utilities, components & mock data
const { useState, useEffect, useRef, useCallback, useMemo } = React;

let _imgCnt = 0;
function Img({ label, bg = '#0c4a6e', stripe = '#083b5a', style = {}, className = '' }) {
  const id = useRef('ip' + (++_imgCnt)).current;
  return (
    <div className={'imgbox ' + className} style={{ background: bg, ...style }}>
      <svg width="100%" height="100%" style={{ position:'absolute',inset:0,top:0,left:0 }} preserveAspectRatio="xMidYMid slice">
        <defs>
          <pattern id={id} width="28" height="28" patternUnits="userSpaceOnUse" patternTransform="rotate(42)">
            <rect width="14" height="28" fill={stripe} opacity="0.28" />
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill={'url(#'+id+')'} />
      </svg>
      <span style={{ position:'relative',zIndex:1,color:'rgba(255,255,255,0.8)',fontFamily:'monospace',fontSize:'10px',textAlign:'center',padding:'10px',fontWeight:600,textShadow:'0 1px 4px rgba(0,0,0,.5)',lineHeight:1.5,maxWidth:'90%' }}>{label}</span>
    </div>
  );
}

function Stars({ score, size = 13 }) {
  return <span className="stars" style={{ fontSize: size }}>{'★'.repeat(Math.floor(score))}{'☆'.repeat(5 - Math.floor(score))}</span>;
}

function RatingRow({ score, count, size = 13 }) {
  return (
    <div style={{ display:'flex', alignItems:'center', gap:5 }}>
      <Stars score={score} size={size} />
      <span style={{ fontSize:size, fontWeight:700 }}>{score.toFixed(1)}</span>
      {count !== undefined && <span style={{ fontSize:12, color:'var(--text3)' }}>({count.toLocaleString()})</span>}
    </div>
  );
}

function NumStepper({ value, onChange, min = 0, max = 99, label, sublabel }) {
  return (
    <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', gap:12 }}>
      <div>
        <div style={{ fontSize:14, fontWeight:600 }}>{label}</div>
        {sublabel && <div style={{ fontSize:12, color:'var(--text3)' }}>{sublabel}</div>}
      </div>
      <div className="nstep">
        <button className="nbtn" onClick={() => onChange(Math.max(min, value-1))}>−</button>
        <span className="nval">{value}</span>
        <button className="nbtn" onClick={() => onChange(Math.min(max, value+1))}>+</button>
      </div>
    </div>
  );
}

function TierBadge({ tier }) {
  const cfg = { 'Đồng':['#a16207','#fef9c3'], 'Bạc':['#6b7280','#f3f4f6'], 'Vàng':['#92400e','#fef3c7'], 'Bạch kim':['#1e3a5f','#e0effe'] };
  const [tc, bg] = cfg[tier] || cfg['Đồng'];
  return <span style={{ display:'inline-flex',alignItems:'center',padding:'2px 10px',borderRadius:9999,fontSize:11,fontWeight:700,background:bg,color:tc,letterSpacing:.5 }}>{tier}</span>;
}

// SVG icons (simple paths only)
const Ico = {
  pin: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C8.1 2 5 5.1 5 9c0 5.3 7 13 7 13s7-7.8 7-13c0-3.9-3.1-7-7-7zm0 9.5a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5z"/></svg>,
  clock: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2zm0 18c-4.4 0-8-3.6-8-8s3.6-8 8-8 8 3.6 8 8-3.6 8-8 8zm.5-13H11v6l5.3 3.2.8-1.3-4.6-2.7V7z"/></svg>,
  person: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/></svg>,
  car: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M18.9 6c-.2-.6-.8-1-1.4-1H6.5c-.7 0-1.2.4-1.4 1L3 12v8c0 .6.4 1 1 1h1c.6 0 1-.4 1-1v-1h12v1c0 .6.4 1 1 1h1c.6 0 1-.4 1-1v-8l-2.1-6zM6.5 16c-.8 0-1.5-.7-1.5-1.5S5.7 13 6.5 13s1.5.7 1.5 1.5-.7 1.5-1.5 1.5zm11 0c-.8 0-1.5-.7-1.5-1.5S16.7 13 17.5 13s1.5.7 1.5 1.5-.7 1.5-1.5 1.5zM5 11l1.5-4.5h11L19 11H5z"/></svg>,
  search: () => <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>,
  bell: () => <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><path d="M18 8a6 6 0 00-12 0c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.7 21a2 2 0 01-3.5 0"/></svg>,
  heart: (filled) => <svg width="16" height="16" viewBox="0 0 24 24" fill={filled?'var(--danger)':'none'} stroke={filled?'var(--danger)':'currentColor'} strokeWidth="2"><path d="M20.8 4.6a5.5 5.5 0 00-7.8 0L12 5.7l-1-1.1a5.5 5.5 0 00-7.8 7.8L12 21.2l8.8-8.8a5.5 5.5 0 000-7.8z"/></svg>,
  share: () => <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.6" y1="13.5" x2="15.4" y2="17.5"/><line x1="15.4" y1="6.5" x2="8.6" y2="10.5"/></svg>,
  chevR: () => <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><path d="M9 18l6-6-6-6"/></svg>,
  chevL: () => <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><path d="M15 18l-6-6 6-6"/></svg>,
  check: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><polyline points="20,6 9,17 4,12"/></svg>,
  x: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>,
  filter: () => <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><line x1="4" y1="6" x2="20" y2="6"/><line x1="8" y1="12" x2="16" y2="12"/><line x1="11" y1="18" x2="13" y2="18"/></svg>,
  globe: () => <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>,
  plane: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M21 16v-2l-8-5V3.5c0-.83-.67-1.5-1.5-1.5S10 2.67 10 3.5V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"/></svg>,
  star2: () => <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>,
  chat: () => <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor"><path d="M20 2H4a2 2 0 00-2 2v18l4-4h14a2 2 0 002-2V4a2 2 0 00-2-2zm-2 10H6v-2h12v2zm0-3H6V7h12v2z"/></svg>,
  gift: () => <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M20 12v10H4V12H2v-2h20v2h-2zm-2 0H6v8h12v-8zm2-5H4V5h3.52A3 3 0 0112 2a3 3 0 014.48 3H20v2zM12 4a1 1 0 00-1 1h2a1 1 0 00-1-1z"/></svg>,
  map: () => <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M20.5 3l-.16.03L15 5.1 9 3 3.36 4.9c-.21.07-.36.25-.36.48V20.5c0 .28.22.5.5.5l.16-.03L9 18.9l6 2.1 5.64-1.9c.21-.07.36-.25.36-.48V3.5c0-.28-.22-.5-.5-.5zM15 19l-6-2.11V5l6 2.11V19z"/></svg>,
  percent: () => <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><path d="M18.4 2.6L4.6 18.4l1.4 1.4L20.2 4zM7 8a1 1 0 100-2 1 1 0 000 2zm0 1a2 2 0 110-4 2 2 0 010 4zm10 7a1 1 0 100-2 1 1 0 000 2zm0 1a2 2 0 110-4 2 2 0 010 4z"/></svg>,
};

function formatPrice(n) {
  if (n >= 1000000) return (n/1000000).toFixed(1).replace('.0','') + ' triệu đ';
  return n.toLocaleString('vi-VN') + 'đ';
}
function formatPriceFull(n) {
  return n.toLocaleString('vi-VN') + 'đ';
}

// ── Mock Data ──────────────────────────────────────────────────
const DESTINATIONS = [
  { id:1, name:'Hạ Long', province:'Quảng Ninh', count:48, bg:'#0c4a6e', stripe:'#083b5a', label:'Vịnh Hạ Long' },
  { id:2, name:'Sapa', province:'Lào Cai', count:35, bg:'#2d6a4f', stripe:'#1b4332', label:'Ruộng bậc thang Sapa' },
  { id:3, name:'Hội An', province:'Quảng Nam', count:62, bg:'#92400e', stripe:'#6b2f09', label:'Phố cổ Hội An' },
  { id:4, name:'Phú Quốc', province:'Kiên Giang', count:74, bg:'#0369a1', stripe:'#024d79', label:'Biển Phú Quốc' },
  { id:5, name:'Đà Lạt', province:'Lâm Đồng', count:29, bg:'#5b21b6', stripe:'#3d1580', label:'Thành phố hoa Đà Lạt' },
  { id:6, name:'Nha Trang', province:'Khánh Hòa', count:56, bg:'#0e7490', stripe:'#0a5570', label:'Biển Nha Trang' },
  { id:7, name:'Đà Nẵng', province:'Đà Nẵng', count:81, bg:'#1d4ed8', stripe:'#1440b0', label:'Cầu Rồng Đà Nẵng' },
  { id:8, name:'Hà Giang', province:'Hà Giang', count:22, bg:'#3f6212', stripe:'#2d4a0d', label:'Đèo Mã Pí Lèng' },
];

const TOURS = [
  { id:1, name:'Hà Nội — Sapa — Fansipan 3N2Đ', duration:'3N2Đ', type:'Tour miền núi', departure:'Hà Nội', dest:'Sapa, Lào Cai', price:2800000, oldPrice:3200000, rating:4.8, reviews:156, seats:5, badge:'Hot', badgeV:'accent', bg:'#2d6a4f', stripe:'#1b4332', label:'Ruộng bậc thang & Fansipan', transport:'Limousine' },
  { id:2, name:'Hạ Long Bay Luxury Cruise 2N1Đ', duration:'2N1Đ', type:'Tour biển đảo', departure:'Hà Nội', dest:'Hạ Long, Quảng Ninh', price:3500000, oldPrice:null, rating:4.9, reviews:312, seats:12, badge:'Best Seller', badgeV:'primary', bg:'#0c4a6e', stripe:'#083b5a', label:'Du thuyền vịnh Hạ Long', transport:'Du thuyền' },
  { id:3, name:'Đà Nẵng — Hội An — Bà Nà Hills 4N3Đ', duration:'4N3Đ', type:'Tour thành phố', departure:'TP.HCM', dest:'Đà Nẵng, Quảng Nam', price:4200000, oldPrice:5000000, rating:4.7, reviews:245, seats:20, badge:'Sale 16%', badgeV:'success', bg:'#92400e', stripe:'#6b2f09', label:'Cầu Vàng Bà Nà Hills', transport:'Máy bay' },
  { id:4, name:'Phú Quốc Thiên Đường Biển 3N2Đ', duration:'3N2Đ', type:'Tour nghỉ dưỡng', departure:'TP.HCM', dest:'Phú Quốc, Kiên Giang', price:5800000, oldPrice:6500000, rating:4.9, reviews:428, seats:3, badge:'Còn 3 chỗ', badgeV:'danger', bg:'#0369a1', stripe:'#024d79', label:'Bãi Sao Phú Quốc', transport:'Máy bay' },
  { id:5, name:'Đà Lạt Mộng Mơ — Thành phố Hoa 3N2Đ', duration:'3N2Đ', type:'Tour văn hóa', departure:'TP.HCM', dest:'Đà Lạt, Lâm Đồng', price:1900000, oldPrice:2200000, rating:4.6, reviews:189, seats:15, badge:null, badgeV:null, bg:'#5b21b6', stripe:'#3d1580', label:'Hồ Xuân Hương Đà Lạt', transport:'Xe khách' },
  { id:6, name:'Thác Bản Giốc — Hang Pắc Bó 3N2Đ', duration:'3N2Đ', type:'Tour mạo hiểm', departure:'Hà Nội', dest:'Cao Bằng', price:2600000, oldPrice:2900000, rating:4.7, reviews:98, seats:20, badge:'Mới', badgeV:'primary', bg:'#3f6212', stripe:'#2d4a0d', label:'Thác Bản Giốc hùng vĩ', transport:'Xe khách' },
  { id:7, name:'Nha Trang — Vinpearl Island 4N3Đ', duration:'4N3Đ', type:'Tour nghỉ dưỡng', departure:'Hà Nội', dest:'Nha Trang, Khánh Hòa', price:4800000, oldPrice:5400000, rating:4.8, reviews:267, seats:10, badge:'Nổi bật', badgeV:'primary', bg:'#0e7490', stripe:'#0a5570', label:'Hòn Tằm Nha Trang', transport:'Máy bay' },
  { id:8, name:'Hà Giang Loop — Địa đầu Tổ Quốc 4N3Đ', duration:'4N3Đ', type:'Tour mạo hiểm', departure:'Hà Nội', dest:'Hà Giang', price:3200000, oldPrice:3500000, rating:4.9, reviews:134, seats:4, badge:'Top Pick', badgeV:'gold', bg:'#3f6212', stripe:'#2d4a0d', label:'Cung đường Hà Giang', transport:'Xe máy / Ô tô' },
];

const HOTELS = [
  { id:1, name:'InterContinental Danang Sun Peninsula', stars:5, type:'Resort', loc:'Đà Nẵng', price:3500000, rating:4.9, reviews:892, bg:'#5b21b6', stripe:'#3d1580', label:'Resort 5 sao Đà Nẵng', amenities:['Hồ bơi','Spa','Nhà hàng','Bãi biển'] },
  { id:2, name:'JW Marriott Phu Quoc Emerald Bay', stars:5, type:'Resort', loc:'Phú Quốc', price:4200000, rating:4.9, reviews:654, bg:'#0369a1', stripe:'#024d79', label:'Resort biển Phú Quốc', amenities:['Hồ bơi vô cực','Spa','Golf','Ăn sáng'] },
  { id:3, name:'Mường Thanh Luxury Hà Nội Centre', stars:4, type:'Hotel', loc:'Hà Nội', price:1200000, rating:4.6, reviews:1240, bg:'#0c4a6e', stripe:'#083b5a', label:'KS trung tâm Hà Nội', amenities:['Hồ bơi','Phòng gym','Ăn sáng','Wifi'] },
  { id:4, name:'Fusion Suites Da Nang Beach', stars:4, type:'Hotel', loc:'Đà Nẵng', price:1800000, rating:4.7, reviews:534, bg:'#92400e', stripe:'#6b2f09', label:'KS bãi biển Mỹ Khê', amenities:['Bãi biển riêng','Spa','Ăn sáng','Wifi'] },
  { id:5, name:'La Siesta Premium Hoi An', stars:4, type:'Boutique Hotel', loc:'Hội An', price:980000, rating:4.8, reviews:892, bg:'#854d0e', stripe:'#6b3c0a', label:'Boutique Hotel phố cổ', amenities:['Hồ bơi','Spa','Ăn sáng','City view'] },
  { id:6, name:'The Anam Cam Ranh Nha Trang', stars:5, type:'Resort', loc:'Nha Trang', price:3100000, rating:4.8, reviews:423, bg:'#0e7490', stripe:'#0a5570', label:'Resort Cam Ranh sang trọng', amenities:['Bãi biển','Hồ bơi','Spa','Golf'] },
  { id:7, name:'Terracotta Resort & Spa Dalat', stars:4, type:'Resort', loc:'Đà Lạt', price:1400000, rating:4.7, reviews:312, bg:'#5b21b6', stripe:'#3d1580', label:'Resort giữa rừng thông', amenities:['Hồ bơi','Spa','Ăn sáng','Núi view'] },
  { id:8, name:'La Veranda Resort Phu Quoc', stars:5, type:'Resort', loc:'Phú Quốc', price:2800000, rating:4.9, reviews:567, bg:'#2d6a4f', stripe:'#1b4332', label:'Resort thuộc địa Pháp', amenities:['Bãi biển riêng','Hồ bơi','Ăn sáng','Bar'] },
];

const COMBOS = [
  { id:1, name:'Combo Phú Quốc 3N2Đ — Vietjet + Resort 5 sao', duration:'3N2Đ', dest:'Phú Quốc', oldPrice:17200000, price:12800000, disc:26, from:'Hà Nội / TP.HCM', includes:['Vé máy bay khứ hồi','Resort 5 sao','Ăn sáng','Tour đảo'], bg:'#0369a1', stripe:'#024d79', label:'Combo Phú Quốc trọn gói' },
  { id:2, name:'Combo Đà Nẵng 4N3Đ — Vietnam Airlines + KS 4 sao', duration:'4N3Đ', dest:'Đà Nẵng', oldPrice:10500000, price:8400000, disc:20, from:'Hà Nội / TP.HCM', includes:['Vé máy bay khứ hồi','KS 4 sao bãi biển','Ăn sáng','Tour Bà Nà Hills'], bg:'#1d4ed8', stripe:'#1440b0', label:'Combo Đà Nẵng bay & ở' },
  { id:3, name:'Combo Sapa Mây Trắng 3N2Đ — Limousine + KS 4 sao', duration:'3N2Đ', dest:'Sapa', oldPrice:7300000, price:6200000, disc:15, from:'Hà Nội', includes:['Xe limousine khứ hồi','KS 4 sao view núi','Ăn sáng','Trekking'], bg:'#2d6a4f', stripe:'#1b4332', label:'Combo Sapa trọn gói' },
];

const BLOGS = [
  { id:1, title:'10 địa điểm không thể bỏ qua khi đến Sapa', excerpt:'Sapa không chỉ có ruộng bậc thang, khám phá những góc khuất ít người biết của thành phố trong mây.', cat:'Cẩm nang', time:'5 phút đọc', date:'10/05/2026', bg:'#2d6a4f', stripe:'#1b4332', label:'Bản Cát Cát Sapa' },
  { id:2, title:'Phú Quốc mùa nào đẹp nhất? Bí kíp đặt phòng giá tốt', excerpt:'Tháng nào nên đến Phú Quốc, khách sạn nào đáng tiền và những điều cần tránh khi du lịch đảo ngọc.', cat:'Mẹo du lịch', time:'7 phút đọc', date:'08/05/2026', bg:'#0369a1', stripe:'#024d79', label:'Hoàng hôn Phú Quốc' },
  { id:3, title:'Hành trình Hà Giang Loop 4N3Đ tự túc từ A đến Z', excerpt:'Chia sẻ kinh nghiệm phượt Hà Giang: phương tiện, lưu trú, ăn uống và những lưu ý an toàn cần biết.', cat:'Phượt ký', time:'10 phút đọc', date:'05/05/2026', bg:'#3f6212', stripe:'#2d4a0d', label:'Mã Pí Lèng Hà Giang' },
];

const USER = { name:'Nguyễn Minh Tuấn', email:'tuannm@email.com', phone:'0901 234 567', tier:'Vàng', points:12480, totalSpend:28500000, joinDate:'01/03/2024' };

const BOOKINGS = [
  { id:'BK260501', product:'Hạ Long Bay Luxury Cruise 2N1Đ', type:'tour', useDate:'15/06/2026', bookDate:'01/05/2026', persons:'2 người lớn', total:7000000, status:'confirmed', statusLabel:'Đã xác nhận', agent:'Halong Star Travel', bg:'#0c4a6e', stripe:'#083b5a', label:'Du thuyền Hạ Long' },
  { id:'BK260418', product:'Mường Thanh Luxury Hà Nội', type:'hotel', useDate:'20/05/2026', bookDate:'18/04/2026', persons:'1 phòng · 2 đêm', total:2400000, status:'paid', statusLabel:'Đã thanh toán', agent:'Mường Thanh Group', bg:'#0c4a6e', stripe:'#083b5a', label:'KS Hà Nội' },
  { id:'BK260310', product:'Đà Nẵng — Hội An — Bà Nà Hills 4N3Đ', type:'tour', useDate:'01/04/2026', bookDate:'10/03/2026', persons:'2 NL · 1 trẻ em', total:12600000, status:'completed', statusLabel:'Hoàn tất', agent:'VietExplore Tours', bg:'#92400e', stripe:'#6b2f09', label:'Bà Nà Hills' },
  { id:'BK260201', product:'Combo Phú Quốc 3N2Đ — Vietjet + Resort 5 sao', type:'combo', useDate:'10/03/2026', bookDate:'01/02/2026', persons:'2 người lớn', total:25600000, status:'completed', statusLabel:'Hoàn tất', agent:'Paradise Travel', bg:'#0369a1', stripe:'#024d79', label:'Resort Phú Quốc' },
  { id:'BK260102', product:'Sapa Mộng Mơ 3N2Đ', type:'tour', useDate:'20/01/2026', bookDate:'02/01/2026', persons:'2 người lớn', total:5600000, status:'cancelled', statusLabel:'Đã hủy', agent:'Mountain Travel', bg:'#2d6a4f', stripe:'#1b4332', label:'Sapa' },
];

const DEPARTURES = [
  { id:1, date:'07/06/2026', dayOfWeek:'Thứ 7', slots:30, booked:25, priceAdult:2800000, priceChild:1960000 },
  { id:2, date:'14/06/2026', dayOfWeek:'Thứ 7', slots:30, booked:18, priceAdult:2800000, priceChild:1960000 },
  { id:3, date:'21/06/2026', dayOfWeek:'Thứ 7', slots:30, booked:30, priceAdult:2800000, priceChild:1960000 },
  { id:4, date:'28/06/2026', dayOfWeek:'Thứ 7', slots:30, booked:8, priceAdult:3200000, priceChild:2240000 },
  { id:5, date:'05/07/2026', dayOfWeek:'Thứ 7', slots:30, booked:12, priceAdult:3200000, priceChild:2240000 },
];

Object.assign(window, {
  useState, useEffect, useRef, useCallback, useMemo,
  Img, Stars, RatingRow, NumStepper, TierBadge, Ico,
  formatPrice, formatPriceFull,
  DESTINATIONS, TOURS, HOTELS, COMBOS, BLOGS, USER, BOOKINGS, DEPARTURES,
});
