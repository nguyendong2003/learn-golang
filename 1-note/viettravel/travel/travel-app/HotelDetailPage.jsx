// HotelDetailPage — Gallery, room types, amenities, booking sidebar
function HotelDetailPage({ navigate, params }) {
  const hotel = params?.hotel || HOTELS[0];
  const [activeTab, setActiveTab] = useState('rooms');
  const [checkIn, setCheckIn] = useState('');
  const [checkOut, setCheckOut] = useState('');
  const [adults, setAdults] = useState(2);
  const [children, setChildren] = useState(0);
  const [selRoom, setSelRoom] = useState(null);
  const [savedHotel, setSavedHotel] = useState(false);
  const [voucher, setVoucher] = useState('');
  const [vOk, setVOk] = useState(false);

  const nights = (checkIn && checkOut)
    ? Math.max(1, Math.round((new Date(checkOut) - new Date(checkIn)) / 86400000))
    : 2;
  const roomPrice = selRoom?.price || hotel.price;
  const subtotal = roomPrice * nights;
  const disc = vOk ? Math.round(subtotal * 0.1) : 0;
  const total = subtotal - disc;

  const roomTypes = [
    { id:1, name:'Phòng Standard', size:28, beds:'1 giường đôi', maxGuests:2, view:'City view', price:Math.round(hotel.price*0.68), amenities:['Smart TV','Wifi','Điều hòa','Minibar'], avail:6, bg:'#475569', stripe:'#334155' },
    { id:2, name:'Phòng Deluxe', size:38, beds:'1 giường King', maxGuests:2, view:['Đà Nẵng','Nha Trang','Phú Quốc'].includes(hotel.loc)?'Sea view':'Garden view', price:hotel.price, amenities:['Smart TV','Wifi','Điều hòa','Minibar','Ban công','Bồn tắm'], avail:3, bg:hotel.bg, stripe:hotel.stripe },
    { id:3, name:'Suite Cao cấp', size:60, beds:'1 King + Sofa bed', maxGuests:3, view:'Panoramic view', price:Math.round(hotel.price*1.9), amenities:['Smart TV','Wifi','Điều hòa','Minibar','Ban công','Bồn tắm','Phòng khách'], avail:1, bg:'#1e3a5f', stripe:'#162d4a' },
  ];

  const tabs = [
    { id:'overview', label:'Tổng quan' },
    { id:'rooms', label:'Loại phòng' },
    { id:'amenities', label:'Tiện nghi' },
    { id:'policy', label:'Chính sách' },
    { id:'reviews', label:'Đánh giá' },
  ];

  const imgColors = [[hotel.bg, hotel.stripe], ['#0e7490','#0a5570'], ['#2d6a4f','#1b4332'], ['#5b21b6','#3d1580']];
  const imgLabels = [hotel.label, 'Hồ bơi & Khu nghỉ dưỡng', 'Nhà hàng & Bar', 'Khu vực thư giãn'];

  return (
    <div className="page-anim" style={{ paddingTop:68 }}>
      {/* Breadcrumb */}
      <div style={{ borderBottom:'1px solid var(--border-lt)', background:'var(--surface)' }}>
        <div className="container" style={{ padding:'12px 32px', display:'flex', gap:6, fontSize:13, color:'var(--text3)', alignItems:'center' }}>
          <button onClick={() => navigate('home')} style={{ background:'none',border:'none',cursor:'pointer',color:'var(--primary)',fontSize:13,padding:0,fontWeight:500 }}>Trang chủ</button>
          <Ico.chevR />
          <button onClick={() => navigate('search',{type:'hotel'})} style={{ background:'none',border:'none',cursor:'pointer',color:'var(--primary)',fontSize:13,padding:0,fontWeight:500 }}>Khách sạn</button>
          <Ico.chevR />
          <span style={{ color:'var(--text)', fontWeight:500 }}>{hotel.loc}</span>
        </div>
      </div>

      <div className="container" style={{ paddingTop:28, paddingBottom:48 }}>
        {/* Title row */}
        <div style={{ marginBottom:16 }}>
          <div style={{ display:'flex', gap:8, marginBottom:10, flexWrap:'wrap', alignItems:'center' }}>
            <span style={{ color:'var(--gold)', fontSize:20, letterSpacing:1 }}>{'★'.repeat(hotel.stars)}</span>
            <span className="badge badge-neutral">{hotel.type}</span>
            <span className="badge badge-neutral" style={{ display:'flex', alignItems:'center', gap:3 }}><Ico.pin />{hotel.loc}</span>
          </div>
          <h1 style={{ fontSize:28, fontWeight:900, letterSpacing:'-.02em', lineHeight:1.2, marginBottom:10 }}>{hotel.name}</h1>
          <div style={{ display:'flex', alignItems:'center', gap:16, flexWrap:'wrap' }}>
            <RatingRow score={hotel.rating} count={hotel.reviews} size={14} />
            <div style={{ display:'flex', gap:6, flexWrap:'wrap' }}>
              {hotel.amenities.map(a => <span key={a} style={{ fontSize:12, padding:'3px 9px', borderRadius:20, background:'var(--border-lt)', color:'var(--text2)', fontWeight:500 }}>{a}</span>)}
            </div>
            <div style={{ display:'flex', gap:8, marginLeft:'auto' }}>
              <button onClick={() => setSavedHotel(!savedHotel)} className="btn btn-ghost btn-sm" style={{ gap:5, color:savedHotel?'var(--danger)':'var(--text2)' }}>
                <Ico.heart filled={savedHotel} />{savedHotel?'Đã lưu':'Lưu'}
              </button>
              <button className="btn btn-ghost btn-sm" style={{ gap:5 }}><Ico.share /> Chia sẻ</button>
            </div>
          </div>
        </div>

        {/* Gallery */}
        <div style={{ display:'grid', gridTemplateColumns:'2fr 1fr', gridTemplateRows:'220px 220px', gap:10, borderRadius:'var(--r-xl)', overflow:'hidden', marginBottom:32 }}>
          <Img bg={imgColors[0][0]} stripe={imgColors[0][1]} label={imgLabels[0]} style={{ gridRow:'1 / 3' }} />
          {[1,2,3].map(i => (
            <div key={i} style={{ position:'relative', overflow:'hidden' }}>
              <Img bg={imgColors[i%4][0]} stripe={imgColors[i%4][1]} label={imgLabels[i%4]} style={{ width:'100%', height:'100%' }} />
              {i===3 && <div style={{ position:'absolute', inset:0, background:'rgba(0,0,0,.45)', display:'flex', alignItems:'center', justifyContent:'center', color:'#fff', fontWeight:700, fontSize:15 }}>+8 ảnh</div>}
            </div>
          ))}
        </div>

        {/* Main layout */}
        <div style={{ display:'grid', gridTemplateColumns:'1fr 348px', gap:28, alignItems:'flex-start' }}>

          {/* Content */}
          <div>
            <div style={{ display:'flex', gap:0, borderBottom:'2px solid var(--border-lt)', marginBottom:28 }}>
              {tabs.map(t => (
                <button key={t.id} onClick={() => setActiveTab(t.id)}
                  style={{ padding:'10px 20px', border:'none', cursor:'pointer', fontWeight:700, fontSize:14, background:'transparent', transition:'all .15s', borderBottom: activeTab===t.id ? '3px solid var(--primary)' : '3px solid transparent', color: activeTab===t.id ? 'var(--primary)' : 'var(--text2)', marginBottom:'-2px' }}>
                  {t.label}
                </button>
              ))}
            </div>
            {activeTab === 'overview' && <HotelOverviewTab hotel={hotel} />}
            {activeTab === 'rooms' && <HotelRoomsTab roomTypes={roomTypes} selRoom={selRoom} setSelRoom={setSelRoom} />}
            {activeTab === 'amenities' && <HotelAmenitiesTab hotel={hotel} />}
            {activeTab === 'policy' && <HotelPolicyTab />}
            {activeTab === 'reviews' && <ReviewsTab tour={{ rating:hotel.rating, reviews:hotel.reviews }} />}
          </div>

          {/* Booking Sidebar */}
          <div style={{ position:'sticky', top:84 }}>
            <div className="card-flat" style={{ borderRadius:'var(--r-xl)', padding:'24px', border:'2px solid var(--border)' }}>
              <div style={{ marginBottom:16 }}>
                <div style={{ fontSize:13, color:'var(--text3)' }}>Giá chỉ từ</div>
                <div style={{ fontSize:26, fontWeight:900, color:'var(--accent)' }}>{formatPriceFull(hotel.price)}<span style={{ fontSize:14, color:'var(--text2)', fontWeight:500 }}>/đêm</span></div>
                <div style={{ fontSize:12, color:'var(--success)', fontWeight:600, display:'flex', gap:4, alignItems:'center', marginTop:4 }}><Ico.check /> Đảm bảo giá tốt nhất</div>
              </div>
              <hr style={{ marginBottom:16 }} />

              {/* Date range */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:10 }}>Ngày lưu trú</div>
                <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:0, border:'1.5px solid var(--border)', borderRadius:'var(--r-md)', overflow:'hidden' }}>
                  {[['Nhận phòng','date',checkIn,setCheckIn],['Trả phòng','date',checkOut,setCheckOut]].map(([lbl,type,val,setter],i) => (
                    <div key={lbl} style={{ padding:'10px 14px', borderRight: i===0 ? '1px solid var(--border-lt)' : 'none' }}>
                      <div style={{ fontSize:10, fontWeight:700, color:'var(--text3)', textTransform:'uppercase', letterSpacing:.5, marginBottom:4 }}>{lbl}</div>
                      <input type={type} value={val} onChange={e => setter(e.target.value)} style={{ border:'none', outline:'none', fontSize:13, fontWeight:600, color:'var(--text)', width:'100%', cursor:'pointer', background:'transparent' }} />
                    </div>
                  ))}
                </div>
                {checkIn && checkOut && <div style={{ fontSize:12, color:'var(--primary)', fontWeight:700, marginTop:6, textAlign:'center', background:'var(--primary-xlt)', borderRadius:'var(--r-sm)', padding:'4px 0' }}>{nights} đêm lưu trú</div>}
              </div>

              {/* Guests */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:10 }}>Số khách</div>
                <div style={{ display:'flex', flexDirection:'column', gap:10, padding:'14px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
                  <NumStepper label="Người lớn" sublabel="≥ 12 tuổi" value={adults} onChange={setAdults} min={1} max={8} />
                  <NumStepper label="Trẻ em" sublabel="< 12 tuổi" value={children} onChange={setChildren} min={0} max={4} />
                </div>
              </div>

              {/* Selected room display */}
              {selRoom && (
                <div style={{ marginBottom:16, padding:'12px 14px', background:'var(--primary-xlt)', borderRadius:'var(--r-md)', border:'1.5px solid var(--primary-lt)' }}>
                  <div style={{ fontSize:11, color:'var(--text3)', marginBottom:2 }}>Phòng đã chọn</div>
                  <div style={{ fontWeight:700, fontSize:14, color:'var(--primary)' }}>{selRoom.name}</div>
                  <div style={{ fontSize:13, color:'var(--text2)' }}>{formatPriceFull(selRoom.price)}/đêm · Tối đa {selRoom.maxGuests} khách</div>
                </div>
              )}

              {/* Voucher */}
              <div style={{ marginBottom:16 }}>
                <div style={{ display:'flex', gap:8 }}>
                  <input className="input" placeholder="Mã giảm giá..." value={voucher} onChange={e => setVoucher(e.target.value.toUpperCase())} style={{ flex:1, fontSize:13 }} />
                  <button onClick={() => voucher && setVOk(true)} className="btn btn-outline btn-sm">Áp dụng</button>
                </div>
                {vOk && <div style={{ marginTop:6, fontSize:12, color:'var(--success)', fontWeight:600, display:'flex', gap:4, alignItems:'center' }}><Ico.check /> Giảm 10% thành công</div>}
              </div>

              {/* Price breakdown */}
              <div style={{ background:'var(--bg)', borderRadius:'var(--r-md)', padding:'14px', marginBottom:16 }}>
                <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--text2)', marginBottom:6 }}>
                  <span>{formatPriceFull(roomPrice)} × {nights} đêm</span>
                  <span>{formatPriceFull(subtotal)}</span>
                </div>
                {disc > 0 && <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--success)', marginBottom:6 }}>
                  <span>Voucher −10%</span><span>−{formatPriceFull(disc)}</span>
                </div>}
                <hr style={{ margin:'8px 0' }} />
                <div style={{ display:'flex', justifyContent:'space-between', fontWeight:800, fontSize:16 }}>
                  <span>Tổng cộng</span>
                  <span style={{ color:'var(--accent)' }}>{formatPriceFull(total)}</span>
                </div>
              </div>

              <button onClick={() => navigate('booking', { hotel, adults, children, nights, room:selRoom, type:'hotel' })}
                className="btn btn-accent btn-lg" style={{ width:'100%', borderRadius:'var(--r-md)', justifyContent:'center' }}>
                Đặt phòng — {formatPrice(total)}
              </button>
              <button className="btn btn-ghost" style={{ width:'100%', marginTop:8, justifyContent:'center', fontSize:13 }}>
                Liên hệ khách sạn trực tiếp
              </button>

              <div style={{ marginTop:14, display:'flex', flexDirection:'column', gap:6 }}>
                {['Miễn phí hủy trước 48h','Không cần thẻ tín dụng','Xác nhận ngay lập tức'].map(t => (
                  <div key={t} style={{ fontSize:12, color:'var(--success)', display:'flex', gap:5, alignItems:'center' }}><Ico.check />{t}</div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function HotelOverviewTab({ hotel }) {
  return (
    <div>
      <div style={{ marginBottom:24 }}>
        <h3 style={{ fontWeight:800, fontSize:18, marginBottom:12 }}>Giới thiệu</h3>
        <p style={{ fontSize:14, lineHeight:1.8, color:'var(--text2)' }}>
          {hotel.name} tọa lạc tại vị trí đắc địa ở {hotel.loc}, là lựa chọn nghỉ dưỡng cao cấp hàng đầu khu vực. Với kiến trúc độc đáo, kết hợp hài hòa giữa nét truyền thống Việt Nam và tiêu chuẩn quốc tế, khách sạn mang đến những trải nghiệm khó quên cho mọi du khách.
        </p>
        <p style={{ fontSize:14, lineHeight:1.8, color:'var(--text2)', marginTop:10 }}>
          Đội ngũ nhân viên chuyên nghiệp, tận tâm, thành thạo đa ngôn ngữ sẵn sàng phục vụ 24/7. Dù bạn đến đây cho kỳ nghỉ gia đình, tuần trăng mật hay chuyến công tác, chúng tôi đảm bảo sự hoàn hảo trong từng chi tiết.
        </p>
      </div>
      <div style={{ marginBottom:24 }}>
        <h3 style={{ fontWeight:800, fontSize:18, marginBottom:14 }}>Thông tin cơ bản</h3>
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:10 }}>
          {[['Check-in','14:00 (Sớm nhất 08:00, phụ thu)'],['Check-out','12:00 (Muộn nhất 18:00, phụ thu)'],['Hạng sao','★'.repeat(hotel.stars)],['Loại hình',hotel.type],['Địa chỉ',hotel.loc + ', Việt Nam'],['Ngôn ngữ','Tiếng Việt · Anh · Trung']].map(([k,v]) => (
            <div key={k} style={{ padding:'12px 16px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
              <div style={{ fontSize:11, color:'var(--text3)', fontWeight:700, textTransform:'uppercase', letterSpacing:.5, marginBottom:3 }}>{k}</div>
              <div style={{ fontSize:14, fontWeight:600, color:'var(--text)' }}>{v}</div>
            </div>
          ))}
        </div>
      </div>
      {/* Map placeholder */}
      <div style={{ borderRadius:'var(--r-lg)', overflow:'hidden', border:'1px solid var(--border)', height:200, position:'relative', display:'flex', alignItems:'center', justifyContent:'center' }}>
        <div style={{ position:'absolute', inset:0, background:'linear-gradient(135deg, #e8f4e8 0%, #d4e8c2 100%)' }} />
        <svg width="100%" height="100%" viewBox="0 0 600 200" style={{ position:'absolute', opacity:.5 }}>
          <rect x="0" y="80" width="600" height="30" fill="#c9ddb0"/>
          <rect x="120" y="30" width="150" height="120" fill="#b8d4a8" opacity=".6"/>
          <circle cx="280" cy="95" r="12" fill="var(--accent)" stroke="#fff" strokeWidth="3"/>
        </svg>
        <div style={{ position:'relative', background:'rgba(255,255,255,.9)', borderRadius:'var(--r-md)', padding:'8px 18px', fontSize:13, fontWeight:600, color:'var(--text2)', display:'flex', gap:6, alignItems:'center' }}>
          <Ico.map /> Xem trên Google Maps
        </div>
      </div>
    </div>
  );
}

function HotelRoomsTab({ roomTypes, selRoom, setSelRoom }) {
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:16 }}>
      {roomTypes.map(r => (
        <div key={r.id} onClick={() => setSelRoom(selRoom?.id===r.id ? null : r)}
          style={{ border: selRoom?.id===r.id ? '2px solid var(--primary)' : '1.5px solid var(--border)', borderRadius:'var(--r-lg)', overflow:'hidden', cursor:'pointer', transition:'all .2s', background: selRoom?.id===r.id ? 'var(--primary-xlt)' : '#fff' }}>
          <div style={{ display:'flex' }}>
            <Img bg={r.bg} stripe={r.stripe} label={r.name} style={{ width:170, height:130, flexShrink:0 }} />
            <div style={{ flex:1, padding:'16px 20px' }}>
              <div style={{ display:'flex', justifyContent:'space-between', alignItems:'flex-start', marginBottom:8 }}>
                <div>
                  <div style={{ fontWeight:800, fontSize:16, marginBottom:5 }}>{r.name}</div>
                  <div style={{ display:'flex', gap:14, fontSize:13, color:'var(--text2)', flexWrap:'wrap' }}>
                    <span>{r.size}m²</span>
                    <span>{r.beds}</span>
                    <span>{r.view}</span>
                    <span>Tối đa {r.maxGuests} khách</span>
                  </div>
                </div>
                <div style={{ textAlign:'right', flexShrink:0, marginLeft:12 }}>
                  <div style={{ fontSize:20, fontWeight:900, color:'var(--accent)' }}>{formatPrice(r.price)}</div>
                  <div style={{ fontSize:12, color:'var(--text3)' }}>/đêm</div>
                  {r.avail <= 3 && <div style={{ fontSize:11, color:'var(--danger)', fontWeight:700, marginTop:2 }}>Còn {r.avail} phòng!</div>}
                </div>
              </div>
              <div style={{ display:'flex', gap:6, flexWrap:'wrap', marginBottom:12 }}>
                {r.amenities.map(a => <span key={a} style={{ fontSize:11, padding:'2px 8px', borderRadius:20, background:'var(--border-lt)', color:'var(--text2)', fontWeight:500 }}>{a}</span>)}
              </div>
              <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center' }}>
                <div style={{ display:'flex', gap:12 }}>
                  <span style={{ fontSize:12, color:'var(--success)', display:'flex', gap:4, alignItems:'center' }}><Ico.check />Miễn phí hủy</span>
                  <span style={{ fontSize:12, color:'var(--success)', display:'flex', gap:4, alignItems:'center' }}><Ico.check />Bao gồm ăn sáng</span>
                </div>
                <button className={'btn btn-sm ' + (selRoom?.id===r.id ? 'btn-primary' : 'btn-outline')} style={{ minWidth:100 }}>
                  {selRoom?.id===r.id ? '✓ Đã chọn' : 'Chọn phòng'}
                </button>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

function HotelAmenitiesTab({ hotel }) {
  const cats = [
    { name:'Tiện nghi chung', items:['Hồ bơi ngoài trời','Phòng gym hiện đại','Spa & Wellness center','Nhà hàng & Rooftop bar','Bãi đỗ xe miễn phí','Sảnh tiếp tân 24/7','Wifi miễn phí toàn khu','Dịch vụ concierge'] },
    { name:'Tiện nghi phòng', items:['Smart TV 55"','Minibar miễn phí','Điều hòa 2 chiều','Két sắt điện tử','Máy sấy tóc cao cấp','Bồn tắm đứng','Ban công riêng','Dép & áo choàng tắm'] },
    { name:'Dịch vụ', items:['Ăn sáng buffet','Đưa đón sân bay (phí)','Giặt ủi & Là trong ngày','Dịch vụ phòng 24/7','Cho thuê xe & Motor','Tour tham quan địa phương','Trẻ em dưới 5 tuổi miễn phí','Hỗ trợ thú cưng nhỏ'] },
  ];
  return (
    <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr', gap:28 }}>
      {cats.map(cat => (
        <div key={cat.name}>
          <h3 style={{ fontWeight:800, fontSize:15, marginBottom:14, color:'var(--text)' }}>{cat.name}</h3>
          <ul style={{ listStyle:'none', display:'flex', flexDirection:'column', gap:9 }}>
            {cat.items.map(item => (
              <li key={item} style={{ fontSize:13, color:'var(--text2)', display:'flex', gap:8, alignItems:'flex-start', lineHeight:1.5 }}>
                <span style={{ color:'var(--success)', flexShrink:0, marginTop:2 }}><Ico.check /></span>{item}
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
}

function HotelPolicyTab() {
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
      <div>
        <h3 style={{ fontWeight:800, fontSize:16, marginBottom:12 }}>Chính sách hủy phòng</h3>
        <div style={{ borderRadius:'var(--r-md)', overflow:'hidden', border:'1px solid var(--border)' }}>
          {[['Hủy trước 48 giờ','Hoàn 100% — Không phí'],['Hủy trước 24–48 giờ','Hoàn 50% — Phí 1 đêm'],['Hủy trong vòng 24 giờ','Không hoàn tiền — Phí 1 đêm'],['No-show','Tính phí toàn bộ đặt phòng']].map(([k,v],i) => (
            <div key={k} style={{ display:'grid', gridTemplateColumns:'1fr 1fr', padding:'14px 20px', background: i%2===0 ? '#fff' : 'var(--bg)', borderBottom: i<3 ? '1px solid var(--border-lt)' : 'none' }}>
              <span style={{ fontSize:14 }}>{k}</span>
              <span style={{ fontSize:14, fontWeight:600, color: v.includes('100%') ? 'var(--success)' : v.includes('Không') ? 'var(--danger)' : 'var(--warn)' }}>{v}</span>
            </div>
          ))}
        </div>
      </div>
      <div>
        <h3 style={{ fontWeight:800, fontSize:16, marginBottom:12 }}>Quy định khách sạn</h3>
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:10 }}>
          {[['Tuổi nhận phòng','Tối thiểu 18 tuổi (có CCCD/Hộ chiếu)'],['Hút thuốc','Chỉ ở khu vực cho phép'],['Tiệc/Ồn ào','Không sau 22:00'],['Thú cưng','Liên hệ trước, phụ thu 200.000đ/đêm'],['Trẻ em','Dưới 5 tuổi miễn phí, dùng giường cha mẹ'],['Thanh toán','Tiền mặt, chuyển khoản, thẻ ngân hàng']].map(([k,v]) => (
            <div key={k} style={{ padding:'12px 16px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
              <div style={{ fontSize:11, color:'var(--text3)', fontWeight:700, marginBottom:3, textTransform:'uppercase', letterSpacing:.4 }}>{k}</div>
              <div style={{ fontSize:13, color:'var(--text)', lineHeight:1.5 }}>{v}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { HotelDetailPage });
