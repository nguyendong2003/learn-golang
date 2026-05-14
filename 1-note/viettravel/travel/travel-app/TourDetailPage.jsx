// TourDetailPage — Gallery, tabs (overview/itinerary/included/policy/reviews), booking sidebar
function TourDetailPage({ navigate, params }) {
  const tour = params?.tour || TOURS[0];
  const [activeTab, setActiveTab] = useState('overview');
  const [selDep, setSelDep] = useState(null);
  const [adults, setAdults] = useState(2);
  const [children, setChildren] = useState(0);
  const [infants, setInfants] = useState(0);
  const [voucher, setVoucher] = useState('');
  const [vApplied, setVApplied] = useState(false);
  const [savedTour, setSavedTour] = useState(false);
  const [mainImg, setMainImg] = useState(0);

  const dep = DEPARTURES.find(d => d.id === selDep);
  const priceAdult = dep ? dep.priceAdult : tour.price;
  const priceChild = dep ? dep.priceChild : Math.round(tour.price * 0.7);
  const subtotal = adults * priceAdult + children * priceChild;
  const discount = vApplied ? Math.round(subtotal * 0.1) : 0;
  const total = subtotal - discount;

  const tabs = [
    { id:'overview', label:'Tổng quan' },
    { id:'itinerary', label:'Lịch trình' },
    { id:'included', label:'Bao gồm' },
    { id:'policy', label:'Chính sách hủy' },
    { id:'reviews', label:'Đánh giá' },
  ];

  const imgLabels = [tour.label, 'Điểm đến ' + tour.dest.split(',')[0], 'Hoạt động du lịch', 'Ẩm thực địa phương'];
  const imgColors = [[tour.bg, tour.stripe], ['#1d4ed8','#1440b0'], ['#0e7490','#0a5570'], ['#3f6212','#2d4a0d']];

  return (
    <div className="page-anim" style={{ paddingTop:68 }}>
      {/* ── Breadcrumb ── */}
      <div style={{ borderBottom:'1px solid var(--border-lt)', background:'var(--surface)' }}>
        <div className="container" style={{ padding:'12px 32px', display:'flex', gap:6, fontSize:13, color:'var(--text3)', alignItems:'center' }}>
          <button onClick={() => navigate('home')} style={{ background:'none', border:'none', cursor:'pointer', color:'var(--primary)', fontSize:13, padding:0, fontWeight:500 }}>Trang chủ</button>
          <Ico.chevR /><span>Tour</span><Ico.chevR />
          <span style={{ color:'var(--text)', fontWeight:500 }}>{tour.type}</span>
        </div>
      </div>

      <div className="container" style={{ paddingTop:28, paddingBottom:48 }}>

        {/* ── Title row ── */}
        <div style={{ marginBottom:16 }}>
          <div style={{ display:'flex', gap:8, marginBottom:10, flexWrap:'wrap' }}>
            <span className="badge badge-primary">{tour.type}</span>
            <span className="badge badge-neutral">{tour.duration}</span>
            <span className="badge badge-neutral" style={{ display:'flex', alignItems:'center', gap:3 }}><Ico.car /> {tour.transport}</span>
          </div>
          <h1 style={{ fontSize:30, fontWeight:900, letterSpacing:'-.02em', lineHeight:1.2, marginBottom:10 }}>{tour.name}</h1>
          <div style={{ display:'flex', alignItems:'center', gap:16, flexWrap:'wrap' }}>
            <RatingRow score={tour.rating} count={tour.reviews} size={14} />
            <div style={{ display:'flex', alignItems:'center', gap:4, fontSize:13, color:'var(--text2)' }}>
              <Ico.pin /> {tour.departure} → {tour.dest}
            </div>
            <div style={{ display:'flex', gap:8, marginLeft:'auto' }}>
              <button onClick={() => setSavedTour(!savedTour)} className="btn btn-ghost btn-sm" style={{ gap:5, color: savedTour?'var(--danger)':'var(--text2)' }}>
                <Ico.heart filled={savedTour} /> {savedTour ? 'Đã lưu' : 'Lưu tour'}
              </button>
              <button className="btn btn-ghost btn-sm" style={{ gap:5 }}><Ico.share /> Chia sẻ</button>
            </div>
          </div>
        </div>

        {/* ── Gallery ── */}
        <div style={{ display:'grid', gridTemplateColumns:'2fr 1fr', gridTemplateRows:'220px 220px', gap:10, borderRadius:'var(--r-xl)', overflow:'hidden', marginBottom:32 }}>
          <Img bg={imgColors[mainImg][0]} stripe={imgColors[mainImg][1]} label={imgLabels[mainImg]} style={{ gridRow:'1 / 3', cursor:'pointer' }} />
          {[1,2,3].map(i => (
            <div key={i} onClick={() => setMainImg(i)} style={{ cursor:'pointer', position:'relative', overflow:'hidden', borderRadius:0 }}>
              <Img bg={imgColors[i%4][0]} stripe={imgColors[i%4][1]} label={imgLabels[i%4]} style={{ width:'100%', height:'100%' }} />
              {i === 3 && <div style={{ position:'absolute', inset:0, background:'rgba(0,0,0,.45)', display:'flex', alignItems:'center', justifyContent:'center', color:'#fff', fontWeight:700, fontSize:15 }}>+12 ảnh</div>}
            </div>
          ))}
        </div>

        {/* ── Main layout ── */}
        <div style={{ display:'grid', gridTemplateColumns:'1fr 348px', gap:28, alignItems:'flex-start' }}>

          {/* ── Content ── */}
          <div>
            {/* Tabs */}
            <div style={{ display:'flex', gap:0, borderBottom:'2px solid var(--border-lt)', marginBottom:28 }}>
              {tabs.map(t => (
                <button key={t.id} onClick={() => setActiveTab(t.id)}
                  style={{ padding:'10px 20px', border:'none', cursor:'pointer', fontWeight:700, fontSize:14, background:'transparent', transition:'all .15s', borderBottom: activeTab===t.id ? '3px solid var(--primary)' : '3px solid transparent', color: activeTab===t.id ? 'var(--primary)' : 'var(--text2)', marginBottom:'-2px' }}>
                  {t.label}
                </button>
              ))}
            </div>

            {/* Tab content */}
            {activeTab === 'overview' && <OverviewTab tour={tour} />}
            {activeTab === 'itinerary' && <ItineraryTab />}
            {activeTab === 'included' && <IncludedTab />}
            {activeTab === 'policy' && <PolicyTab />}
            {activeTab === 'reviews' && <ReviewsTab tour={tour} />}
          </div>

          {/* ── Booking Sidebar ── */}
          <div style={{ position:'sticky', top:84 }}>
            <div className="card-flat" style={{ borderRadius:'var(--r-xl)', padding:'24px', border:'2px solid var(--border)' }}>
              <div style={{ marginBottom:16 }}>
                {tour.oldPrice && <div style={{ fontSize:13, color:'var(--text3)', textDecoration:'line-through' }}>{formatPriceFull(tour.oldPrice)}</div>}
                <div style={{ fontSize:26, fontWeight:900, color:'var(--accent)' }}>{formatPriceFull(priceAdult)}<span style={{ fontSize:14, color:'var(--text2)', fontWeight:500 }}>/người lớn</span></div>
                <div style={{ fontSize:13, color:'var(--text2)' }}>Đã gồm: VAT, hướng dẫn viên, bảo hiểm</div>
              </div>
              <hr style={{ marginBottom:16 }} />

              {/* Chọn đợt khởi hành */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:10 }}>Chọn đợt khởi hành</div>
                <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
                  {DEPARTURES.map(d => {
                    const full = d.booked >= d.slots;
                    const few = d.slots - d.booked <= 5 && !full;
                    return (
                      <button key={d.id} onClick={() => !full && setSelDep(d.id)}
                        disabled={full}
                        style={{ padding:'10px 14px', borderRadius:'var(--r-md)', border: selDep===d.id ? '2px solid var(--primary)' : '1.5px solid var(--border)', background: selDep===d.id ? 'var(--primary-xlt)' : full ? 'var(--border-lt)' : '#fff', cursor: full?'not-allowed':'pointer', textAlign:'left', transition:'all .15s', opacity: full?0.6:1 }}>
                        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center' }}>
                          <div>
                            <span style={{ fontWeight:700, fontSize:14, color: selDep===d.id ? 'var(--primary)' : 'var(--text)' }}>{d.date}</span>
                            <span style={{ fontSize:12, color:'var(--text3)', marginLeft:6 }}>{d.dayOfWeek}</span>
                          </div>
                          {full ? <span className="badge badge-danger">Hết chỗ</span>
                            : few ? <span className="badge badge-warn">Còn {d.slots-d.booked} chỗ</span>
                            : <span className="badge badge-success">Còn chỗ</span>}
                        </div>
                        <div style={{ fontSize:13, color:'var(--text2)', marginTop:3 }}>{formatPriceFull(d.priceAdult)}/người lớn</div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Số hành khách */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:10 }}>Số hành khách</div>
                <div style={{ display:'flex', flexDirection:'column', gap:10, padding:'14px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
                  <NumStepper label="Người lớn" sublabel="≥ 12 tuổi" value={adults} onChange={setAdults} min={1} max={20} />
                  <NumStepper label="Trẻ em" sublabel="5–11 tuổi · 70% giá" value={children} onChange={setChildren} min={0} max={10} />
                  <NumStepper label="Em bé" sublabel="< 5 tuổi · Miễn phí" value={infants} onChange={setInfants} min={0} max={5} />
                </div>
              </div>

              {/* Voucher */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:8 }}>Mã giảm giá</div>
                <div style={{ display:'flex', gap:8 }}>
                  <input className="input" placeholder="Nhập mã voucher..." value={voucher} onChange={e => setVoucher(e.target.value.toUpperCase())} style={{ flex:1, fontSize:13 }} />
                  <button onClick={() => voucher && setVApplied(true)} className="btn btn-outline btn-sm">Áp dụng</button>
                </div>
                {vApplied && <div style={{ marginTop:6, fontSize:12, color:'var(--success)', fontWeight:600, display:'flex', gap:4, alignItems:'center' }}><Ico.check /> Áp dụng thành công — Giảm 10%</div>}
              </div>

              {/* Price breakdown */}
              <div style={{ background:'var(--bg)', borderRadius:'var(--r-md)', padding:'14px', marginBottom:16 }}>
                <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--text2)', marginBottom:6 }}>
                  <span>Người lớn × {adults}</span><span>{formatPriceFull(adults * priceAdult)}</span>
                </div>
                {children > 0 && <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--text2)', marginBottom:6 }}>
                  <span>Trẻ em × {children}</span><span>{formatPriceFull(children * priceChild)}</span>
                </div>}
                {discount > 0 && <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--success)', marginBottom:6 }}>
                  <span>Giảm giá voucher</span><span>−{formatPriceFull(discount)}</span>
                </div>}
                <hr style={{ margin:'8px 0' }} />
                <div style={{ display:'flex', justifyContent:'space-between', fontWeight:800, fontSize:16 }}>
                  <span>Tổng cộng</span><span style={{ color:'var(--accent)' }}>{formatPriceFull(total)}</span>
                </div>
              </div>

              <button onClick={() => navigate('booking', { tour, dep, adults, children })}
                className="btn btn-accent btn-lg" style={{ width:'100%', borderRadius:'var(--r-md)', justifyContent:'center' }}>
                Đặt ngay — {formatPrice(total)}
              </button>
              <button className="btn btn-ghost" style={{ width:'100%', marginTop:8, justifyContent:'center', fontSize:13 }}>
                Liên hệ đại lý tư vấn
              </button>

              {/* Agent card */}
              <div style={{ marginTop:16, padding:'14px', background:'var(--primary-xlt)', borderRadius:'var(--r-md)', display:'flex', gap:12, alignItems:'center' }}>
                <div style={{ width:44, height:44, borderRadius:'50%', background:'var(--primary)', display:'flex', alignItems:'center', justifyContent:'center', color:'#fff', fontWeight:800, fontSize:16, flexShrink:0 }}>HT</div>
                <div style={{ flex:1 }}>
                  <div style={{ fontWeight:700, fontSize:13 }}>Halong Star Travel</div>
                  <div style={{ fontSize:12, color:'var(--text2)', display:'flex', gap:4, alignItems:'center' }}><Stars score={4.8} size={10} /> 4.8 · 312 đánh giá</div>
                  <span className="badge badge-success" style={{ marginTop:4, fontSize:10 }}>Đại lý chứng nhận</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function OverviewTab({ tour }) {
  const highlights = ['Chinh phục đỉnh Fansipan hùng vĩ — nóc nhà Đông Dương 3.147m','Trải nghiệm ruộng bậc thang tuyệt đẹp ở Mường Hoa, Tả Van','Thăm bản làng dân tộc H\'Mông, Dao tại Cát Cát, Tả Phìn','Ngủ homestay truyền thống, thưởng thức ẩm thực địa phương','Cáp treo Fansipan — tuyến cáp treo 3 dây dài nhất Đông Nam Á'];
  return (
    <div>
      <div style={{ marginBottom:24 }}>
        <h3 style={{ fontWeight:800, fontSize:18, marginBottom:12 }}>Điểm nổi bật</h3>
        <ul style={{ listStyle:'none', display:'flex', flexDirection:'column', gap:8 }}>
          {highlights.map((h,i) => (
            <li key={i} style={{ display:'flex', gap:10, alignItems:'flex-start', fontSize:14, lineHeight:1.5 }}>
              <span style={{ color:'var(--success)', flexShrink:0, marginTop:2 }}><Ico.check /></span>
              <span>{h}</span>
            </li>
          ))}
        </ul>
      </div>
      <div style={{ marginBottom:24 }}>
        <h3 style={{ fontWeight:800, fontSize:18, marginBottom:12 }}>Mô tả tour</h3>
        <p style={{ fontSize:14, lineHeight:1.8, color:'var(--text2)' }}>Hành trình Sapa – Fansipan 3 ngày 2 đêm đưa bạn đến với vùng đất núi rừng hùng vĩ của Tây Bắc Việt Nam. Từ những thửa ruộng bậc thang trải dài như dải lụa đến đỉnh Fansipan sương mù bí ẩn — chuyến đi này hứa hẹn những trải nghiệm không thể nào quên.</p>
        <p style={{ fontSize:14, lineHeight:1.8, color:'var(--text2)', marginTop:10 }}>Chương trình tour được thiết kế linh hoạt, kết hợp giữa tham quan, khám phá văn hóa bản địa và nghỉ dưỡng. Đoàn được đồng hành bởi hướng dẫn viên kinh nghiệm, am hiểu văn hóa địa phương.</p>
      </div>
      <MapPlaceholder />
    </div>
  );
}

function MapPlaceholder() {
  return (
    <div style={{ borderRadius:'var(--r-lg)', overflow:'hidden', border:'1px solid var(--border)', marginBottom:8 }}>
      <div style={{ background:'#e8f0e0', height:280, position:'relative', display:'flex', alignItems:'center', justifyContent:'center' }}>
        <div style={{ position:'absolute', inset:0 }}>
          <svg width="100%" height="100%" viewBox="0 0 600 280" preserveAspectRatio="xMidYMid slice">
            <rect width="600" height="280" fill="#e8f0e0"/>
            <rect x="0" y="120" width="600" height="40" fill="#d4e8c2" opacity=".6"/>
            <rect x="140" y="60" width="180" height="140" fill="#c9ddb0" opacity=".4"/>
            <circle cx="200" cy="130" r="8" fill="var(--primary)" stroke="#fff" strokeWidth="3"/>
            <circle cx="380" cy="100" r="8" fill="var(--accent)" stroke="#fff" strokeWidth="3"/>
            <line x1="200" y1="130" x2="380" y2="100" stroke="var(--primary)" strokeWidth="2" strokeDasharray="6,4"/>
          </svg>
        </div>
        <div style={{ position:'relative', zIndex:1, textAlign:'center' }}>
          <div style={{ background:'rgba(255,255,255,.9)', borderRadius:'var(--r-md)', padding:'10px 20px', fontSize:13, fontWeight:600, color:'var(--text2)' }}>
            <Ico.map /> &nbsp; Bản đồ lộ trình tour
          </div>
        </div>
      </div>
    </div>
  );
}

function ItineraryTab() {
  const days = [
    { day:'Ngày 1', title:'Hà Nội → Sapa', meals:['Tối'], hotel:'Sapa House Hotel ⭐⭐⭐⭐', desc:'07:00 Xe limousine đón tại điểm tập kết Hà Nội khởi hành đi Sapa. 12:00 Ăn trưa tại nhà hàng địa phương giữa đường. 16:00 Đến Sapa, nhận phòng khách sạn. 17:30 Tự do khám phá thị trấn Sapa và chợ đêm. 19:30 Ăn tối tại nhà hàng với đặc sản Tây Bắc.' },
    { day:'Ngày 2', title:'Sapa → Trekking Cát Cát → Fansipan', meals:['Sáng','Trưa','Tối'], hotel:'Sapa House Hotel ⭐⭐⭐⭐', desc:'07:00 Ăn sáng tại khách sạn. 08:00 Trekking bản Cát Cát — thăm làng người H\'Mông cổ truyền, ngắm thác nước và cầu treo. 11:00 Bắt cáp treo Fansipan lên đỉnh 3.147m — chụp ảnh lưu niệm. 13:00 Ăn trưa. 15:00 Tham quan Ga Mây trên đỉnh Fansipan. 19:00 Ăn tối — thưởng thức thắng cố, rượu táo mèo.' },
    { day:'Ngày 3', title:'Tả Van → Mường Hoa → Hà Nội', meals:['Sáng'], hotel:'', desc:'07:00 Ăn sáng, trả phòng. 08:00 Trekking bản Tả Van và thung lũng Mường Hoa — ngắm ruộng bậc thang nổi tiếng. 11:00 Tự do mua sắm đặc sản địa phương. 12:00 Xe limousine khởi hành về Hà Nội. 18:00 Về đến Hà Nội, kết thúc chương trình.' },
  ];
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
      {days.map((d, i) => (
        <div key={i} style={{ display:'flex', gap:20 }}>
          <div style={{ flexShrink:0, display:'flex', flexDirection:'column', alignItems:'center' }}>
            <div style={{ width:44, height:44, borderRadius:'50%', background:'var(--primary)', color:'#fff', display:'flex', alignItems:'center', justifyContent:'center', fontWeight:800, fontSize:12, textAlign:'center', lineHeight:1.2 }}>{d.day.replace('Ngày ','N')}</div>
            {i < days.length-1 && <div style={{ width:2, flex:1, background:'var(--border-lt)', margin:'6px 0' }} />}
          </div>
          <div style={{ flex:1, paddingBottom: i < days.length-1 ? 20 : 0 }}>
            <div style={{ fontWeight:800, fontSize:16, marginBottom:6 }}>{d.day}: {d.title}</div>
            <div style={{ display:'flex', gap:6, marginBottom:10, flexWrap:'wrap' }}>
              {d.meals.map(m => <span key={m} style={{ fontSize:12, padding:'2px 9px', borderRadius:20, background:'var(--gold-lt)', color:'var(--gold)', fontWeight:600 }}>🍽 Bữa {m}</span>)}
              {d.hotel && <span style={{ fontSize:12, padding:'2px 9px', borderRadius:20, background:'var(--primary-lt)', color:'var(--primary)', fontWeight:600 }}>🏨 {d.hotel}</span>}
            </div>
            <p style={{ fontSize:14, lineHeight:1.7, color:'var(--text2)' }}>{d.desc}</p>
          </div>
        </div>
      ))}
    </div>
  );
}

function IncludedTab() {
  const inc = ['Xe limousine khứ hồi Hà Nội — Sapa','Khách sạn 4 sao tại Sapa (2 đêm, phòng đôi)','Ăn sáng 2 ngày tại khách sạn','Ăn trưa & tối theo chương trình (4 bữa)','Vé cáp treo Fansipan khứ hồi','Hướng dẫn viên tiếng Việt suốt hành trình','Vé tham quan các điểm trong chương trình','Bảo hiểm du lịch 50 triệu đồng','Nước khoáng trên xe (2 chai/ngày)'];
  const excl = ['Chi phí cá nhân (giặt ủi, điện thoại, mua sắm)','Đồ uống tại nhà hàng','Phụ thu phòng đơn: 350.000đ/đêm','Tip cho hướng dẫn viên và lái xe','Chi phí vượt ngoài chương trình'];
  return (
    <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:28 }}>
      <div>
        <h3 style={{ fontWeight:800, fontSize:16, marginBottom:14, color:'var(--success)', display:'flex', gap:6, alignItems:'center' }}><Ico.check /> Tour bao gồm</h3>
        <ul style={{ listStyle:'none', display:'flex', flexDirection:'column', gap:8 }}>
          {inc.map((item,i) => <li key={i} style={{ display:'flex', gap:8, fontSize:14, color:'var(--text2)', lineHeight:1.5, alignItems:'flex-start' }}><span style={{ color:'var(--success)', flexShrink:0, marginTop:2 }}><Ico.check /></span>{item}</li>)}
        </ul>
      </div>
      <div>
        <h3 style={{ fontWeight:800, fontSize:16, marginBottom:14, color:'var(--danger)', display:'flex', gap:6, alignItems:'center' }}><Ico.x /> Không bao gồm</h3>
        <ul style={{ listStyle:'none', display:'flex', flexDirection:'column', gap:8 }}>
          {excl.map((item,i) => <li key={i} style={{ display:'flex', gap:8, fontSize:14, color:'var(--text2)', lineHeight:1.5, alignItems:'flex-start' }}><span style={{ color:'var(--danger)', flexShrink:0, marginTop:2 }}><Ico.x /></span>{item}</li>)}
        </ul>
      </div>
    </div>
  );
}

function PolicyTab() {
  const rows = [['Hủy trước 15 ngày','Hoàn 100% — Không phí'],['Hủy trước 7–14 ngày','Hoàn 70% — Phí 30%'],['Hủy trước 3–6 ngày','Hoàn 50% — Phí 50%'],['Hủy dưới 3 ngày','Không hoàn tiền'],['Không đến (No-show)','Không hoàn tiền']];
  return (
    <div>
      <h3 style={{ fontWeight:800, fontSize:18, marginBottom:16 }}>Chính sách hủy tour</h3>
      <div style={{ borderRadius:'var(--r-md)', overflow:'hidden', border:'1px solid var(--border)' }}>
        {rows.map(([cond, policy], i) => (
          <div key={i} style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:0, padding:'14px 20px', background: i%2===0 ? '#fff' : 'var(--bg)', borderBottom: i<rows.length-1 ? '1px solid var(--border-lt)' : 'none' }}>
            <span style={{ fontSize:14, fontWeight:500 }}>{cond}</span>
            <span style={{ fontSize:14, color: policy.includes('100%') ? 'var(--success)' : policy.includes('Không') ? 'var(--danger)' : 'var(--warn)', fontWeight:600 }}>{policy}</span>
          </div>
        ))}
      </div>
      <div style={{ marginTop:16, padding:'14px 18px', background:'var(--warn-lt)', borderRadius:'var(--r-md)', fontSize:13, color:'var(--text2)', lineHeight:1.6 }}>
        <strong style={{ color:'var(--warn)' }}>Lưu ý:</strong> Chính sách hủy áp dụng tính theo ngày dương lịch. Yêu cầu hủy phải được gửi bằng email hoặc qua hệ thống. Thời gian hoàn tiền 7–14 ngày làm việc.
      </div>
    </div>
  );
}

function ReviewsTab({ tour }) {
  const reviews = [
    { name:'Trần Văn Minh', avatar:'TM', date:'20/04/2026', score:5, text:'Tour tuyệt vời! Hướng dẫn viên nhiệt tình và chuyên nghiệp. Khách sạn đẹp, sạch sẽ, view nhìn ra ruộng bậc thang cực đẹp. Chắc chắn sẽ quay lại!', imgs: true },
    { name:'Lê Thị Hoa', avatar:'LH', date:'15/04/2026', score:5, text:'Chuyến đi Sapa lần đầu tiên của gia đình mình, mọi thứ đều rất ổn. Xe đúng giờ, khách sạn tốt, đồ ăn ngon. Con bé nhà mình thích lắm. Cáp treo Fansipan hơi chờ lâu nhưng worth it!', imgs: false },
    { name:'Nguyễn Hoàng An', avatar:'HA', date:'10/04/2026', score:4, text:'Tour ổn, nhưng ngày 2 trời mưa nên trekking khá vất vả. Mình nghĩ nên có plan B cho thời tiết xấu. Ngoài ra thì hướng dẫn viên rất thú vị và am hiểu văn hóa địa phương.', imgs: false },
  ];
  const dist = [5,4,3,2,1];
  const counts = [102,38,12,3,1];
  const total = counts.reduce((a,b)=>a+b,0);
  return (
    <div>
      {/* Summary */}
      <div style={{ display:'flex', gap:32, marginBottom:28, padding:'24px', background:'var(--bg)', borderRadius:'var(--r-lg)' }}>
        <div style={{ textAlign:'center' }}>
          <div style={{ fontSize:52, fontWeight:900, color:'var(--text)', lineHeight:1 }}>{tour.rating}</div>
          <Stars score={tour.rating} size={18} />
          <div style={{ fontSize:13, color:'var(--text3)', marginTop:4 }}>{tour.reviews.toLocaleString()} đánh giá</div>
        </div>
        <div style={{ flex:1 }}>
          {dist.map((d,i) => (
            <div key={d} style={{ display:'flex', alignItems:'center', gap:10, marginBottom:6 }}>
              <span style={{ fontSize:13, color:'var(--text2)', width:30, textAlign:'right' }}>{d}★</span>
              <div style={{ flex:1, height:8, background:'var(--border-lt)', borderRadius:4, overflow:'hidden' }}>
                <div style={{ height:'100%', background:'var(--gold)', borderRadius:4, width: (counts[i]/total*100)+'%', transition:'width .4s' }} />
              </div>
              <span style={{ fontSize:12, color:'var(--text3)', width:24 }}>{counts[i]}</span>
            </div>
          ))}
        </div>
      </div>
      {/* Reviews list */}
      <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
        {reviews.map((r,i) => (
          <div key={i} style={{ padding:'20px', background:'var(--surface)', borderRadius:'var(--r-md)', border:'1px solid var(--border-lt)' }}>
            <div style={{ display:'flex', gap:12, alignItems:'flex-start', marginBottom:10 }}>
              <div style={{ width:40, height:40, borderRadius:'50%', background:'var(--primary)', color:'#fff', display:'flex', alignItems:'center', justifyContent:'center', fontSize:14, fontWeight:800, flexShrink:0 }}>{r.avatar}</div>
              <div style={{ flex:1 }}>
                <div style={{ fontWeight:700, fontSize:14 }}>{r.name}</div>
                <div style={{ display:'flex', alignItems:'center', gap:8 }}>
                  <Stars score={r.score} size={12} />
                  <span style={{ fontSize:12, color:'var(--text3)' }}>{r.date}</span>
                </div>
              </div>
            </div>
            <p style={{ fontSize:14, lineHeight:1.7, color:'var(--text2)' }}>{r.text}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

Object.assign(window, { TourDetailPage, PolicyTab, ReviewsTab });
