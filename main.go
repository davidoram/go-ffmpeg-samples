package main

import (
	"log"
	"os"

	"github.com/davidoram/gmf"
)

func fatal(e error) {
	log.Fatal(e)
}

func abort(s string) {
	log.Fatal(s)
}

type FilterContext struct {
	src  *gmf.AVFilter
	sink *gmf.AVFilter

	inputs  *gmf.AVFilterInOut
	outputs *gmf.AVFilterInOut
}

func initFilters() FilterContext {
	src, err := gmf.GetFilter("buffer")
	if err != nil {
		fatal(err)
	}
	sink, err := gmf.GetFilter("buffersink")
	if err != nil {
		fatal(err)
	}

	inputs, err := gmf.NewFilterInOut()
	if err != nil {
		fatal(err)
	}
	outputs, err := gmf.NewFilterInOut()
	if err != nil {
		fatal(err)
	}

	return FilterContext{src: src, sink: sink, inputs: inputs, outputs: outputs}
}

func (this *FilterContext) freeFilters() {
	this.inputs.Free()
	this.outputs.Free()
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: %s file\n", os.Args[0])
	}

	ictx := gmf.NewCtx()
	log.Printf("Opening input '%v'...\n", os.Args[1])
	ictx.OpenInput(os.Args[1])
	log.Println("Retrieving best stream...")
	ist, err := ictx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		fatal(err)
	}

	filters := initFilters()
	defer filters.freeFilters()

	log.Println("Reading packets...")
	var i = 0
	for p := range ictx.GetNewPackets() {
		_ = p
		i++
		gmf.Release(p)
	}
	log.Println("Done", i, ist)

	ictx.CloseInputAndRelease()

	//   if ((ret = init_filters(filter_descr)) < 0)
	//       goto end;
	//   /* read all packets */
	//   while (1) {
	//       if ((ret = av_read_frame(fmt_ctx, &packet)) < 0)
	//           break;
	//       if (packet.stream_index == video_stream_index) {
	//           got_frame = 0;
	//           ret = avcodec_decode_video2(dec_ctx, frame, &got_frame, &packet);
	//           if (ret < 0) {
	//               av_log(NULL, AV_LOG_ERROR, "Error decoding video\n");
	//               break;
	//           }
	//           if (got_frame) {
	//               frame->pts = av_frame_get_best_effort_timestamp(frame);
	//               /* push the decoded frame into the filtergraph */
	//               if (av_buffersrc_add_frame_flags(buffersrc_ctx, frame, AV_BUFFERSRC_FLAG_KEEP_REF) < 0) {
	//                   av_log(NULL, AV_LOG_ERROR, "Error while feeding the filtergraph\n");
	//                   break;
	//               }
	//               /* pull filtered frames from the filtergraph */
	//               while (1) {
	//                   ret = av_buffersink_get_frame(buffersink_ctx, filt_frame);
	//                   if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
	//                       break;
	//                   if (ret < 0)
	//                       goto end;
	//                   display_frame(filt_frame, buffersink_ctx->inputs[0]->time_base);
	//                   av_frame_unref(filt_frame);
	//               }
	//               av_frame_unref(frame);
	//           }
	//       }
	//       av_free_packet(&packet);
	//   }
	// end:
	//   avfilter_graph_free(&filter_graph);
	//   avcodec_close(dec_ctx);
	//   avformat_close_input(&fmt_ctx);
	//   av_frame_free(&frame);
	//   av_frame_free(&filt_frame);
	//   if (ret < 0 && ret != AVERROR_EOF) {
	//       fprintf(stderr, "Error occurred: %s\n", av_err2str(ret));
	//       exit(1);
	//   }
	//   exit(0);
}
