package tool

//func (m *Media) HasStreams() bool {
//  if m.HasMeta() && m.Meta.Streams != nil {
//    return true
//  }
//  return false
//}

//func (m *Media) VideoCodec() string {
//  if m.HasMeta() {
//    for _, stream := range m.Meta.Streams {
//      if stream.CodecType == "video" {
//        return stream.CodecName
//      }
//    }
//  }
//return ""
//}

//func (m *Media) AudioCodec() string {
//  if m.HasMeta() {
//    for _, stream := range m.Meta.Streams {
//      if stream.CodecType == "audio" {
//        return stream.CodecName
//      }
//    }
//  }
//  return ""
//}

//func (m *Media) HasVideo() bool {
//  if m.HasMeta() {
//    for _, stream := range m.Meta.Streams {
//      if stream.CodecType == "video" {
//        return true
//      }
//    }
//  }
//  return false
//}

//func (m *Media) HasAudio() bool {
//  if m.HasMeta() {
//    for _, stream := range m.Meta.Streams {
//      if stream.CodecType == "audio" {
//        return true
//      }
//    }
//  }
//  return false
//}
